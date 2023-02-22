/* Copyright 2022 Zinc Labs Inc. and Contributors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package directory

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/index"
	segment "github.com/blugelabs/bluge_segment_api"
	"github.com/rs/zerolog/log"
	zincConfig "github.com/zinclabs/zincsearch/pkg/config"
)

// GetS3Config returns a bluge config that will store index data in S3
// bucket: the S3 bucket to use
// indexName: the name of the index to use. It will be an s3 prefix (folder)
func GetS3Config(bucket string, indexName string, timeRange ...int64) bluge.Config {
	config := index.DefaultConfigWithDirectory(func() index.Directory {
		return NewS3Directory(bucket, indexName)
	})
	config = config.WithPersisterNapTimeMSec(50)
	if len(timeRange) == 2 {
		if timeRange[0] <= timeRange[1] {
			config = config.WithTimeRange(timeRange[0], timeRange[1])
		}
	}
	return bluge.DefaultConfigWithIndexConfig(config)
}

type S3Directory struct {
	Bucket string
	Prefix string
	Client *s3.Client
}

var endpointResolver = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
	if zincConfig.Global.S3.Url != "" {
		return aws.Endpoint{
			URL:           zincConfig.Global.S3.Url,
			SigningRegion: region,
			Source:        aws.EndpointSourceCustom,
		}, nil
	}
	// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
	return aws.Endpoint{}, &aws.EndpointNotFoundError{}
})

// NewS3Directory creates a new S3Directory instance which can be used to create s3 backed indexes
func NewS3Directory(bucket, prefix string) index.Directory {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(endpointResolver))
	if err != nil {
		log.Print("Error loading AWS config: ", err)
	}
	client := s3.NewFromConfig(cfg)

	directory := &S3Directory{
		Bucket: bucket,
		Prefix: prefix,
		Client: client,
	}

	return directory
}

func (s *S3Directory) fileName(kind string, id uint64) string {
	return fmt.Sprintf("%012x", id) + kind
}

func (s *S3Directory) Setup(readOnly bool) error {
	return nil
}

// List the ids of all the items of the specified kind
// Items are returned in descending order by id
func (s *S3Directory) List(kind string) ([]uint64, error) {
	log.Print("List: s3 ListObjectsV2 call made for: s3://", s.Bucket+"/"+s.Prefix)
	var itemList []uint64

	ctx := context.Background()
	params := s3.ListObjectsV2Input{
		Bucket: &s.Bucket,
		Prefix: &s.Prefix,
	}

	val, err := s.Client.ListObjectsV2(ctx, &params)
	if err != nil {
		log.Print("List: failed to list objects: ", err.Error())
		return nil, err
	}

	for _, obj := range val.Contents {
		if filepath.Ext(*obj.Key) != kind {
			continue
		}

		stringID := filepath.Base(*obj.Key)
		stringID = stringID[:len(stringID)-len(kind)]

		parsedID, err := strconv.ParseUint(stringID, 16, 64)
		if err != nil {
			log.Print("List: failed to parse object id: ", err.Error())
			continue
		}

		itemList = append(itemList, parsedID)

	}

	return itemList, nil
}

// Load the specified item
// Item data is accessible via the returned *segment.Data structure
// A io.Closer is returned, which must be called to release
// resources held by this open item.
// NOTE: care must be taken to handle a possible nil io.Closer
func (s *S3Directory) Load(kind string, id uint64) (*segment.Data, io.Closer, error) {

	ctx := context.Background()

	key := s.Prefix + "/" + s.fileName(kind, id)

	goi := &s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	}

	log.Print("Load: s3 GetObject call made. s3://", s.Bucket, "/", key)

	output, err := s.Client.GetObject(ctx, goi)

	if err != nil {
		log.Print("Load: failed to get object: s3://"+s.Bucket+"/"+key, err.Error())
		return nil, nil, err
	}
	defer func() {
		if err := output.Body.Close(); err != nil {
			log.Print("Load: failed to close download stream: s3://"+s.Bucket+"/"+key, err.Error())
		}
	}()

	data, err := ioutil.ReadAll(output.Body)
	if err != nil {
		log.Print("Load: failed to read object", err.Error())
		return nil, nil, err
	}

	return segment.NewDataBytes(data), nil, nil
}

// Persist a new item with data from the provided WriterTo
// Implementations should monitor the closeCh and return with error
// in the event it is closed before completion.
func (s *S3Directory) Persist(kind string, id uint64, w index.WriterTo, closeCh chan struct{}) error {
	var buf bytes.Buffer
	_, err := w.WriteTo(&buf, closeCh)
	if err != nil {
		log.Print("Persist: failed to write object to buffer: ", err.Error())
		return err
	}

	s3ObjectName := s.fileName(kind, id)

	path := filepath.Join(s.Prefix, s3ObjectName)

	ctx := context.Background()

	params := s3.PutObjectInput{
		Bucket: &s.Bucket,
		Key:    &path,
		Body:   bytes.NewReader(buf.Bytes()),
	}

	output, err := s.Client.PutObject(ctx, &params)

	if err != nil {
		log.Print("Persist: failed to write object: ", err.Error())
		return err
	}

	h := md5.New()
	h.Write(buf.Bytes())

	if output.ETag == nil || *output.ETag != fmt.Sprintf("%x", h.Sum(nil)) {
		log.Print("Warning: s3 object " + s.Bucket + "/" + path + " has mismatched md5 checksum")
	}

	log.Print("Persist: s3 object "+s.Bucket+"/"+path+" written. Its md5 hash is: ", *output.ETag)
	return nil
}

// Remove the specified item
func (s *S3Directory) Remove(kind string, id uint64) error {
	objectToDelete := filepath.Join(s.Prefix, s.fileName(kind, id))
	ctx := context.Background()
	doi := &s3.DeleteObjectInput{
		Bucket: &s.Bucket,
		Key:    &objectToDelete,
	}

	log.Print("Remove: s3 DeleteObject call made s3://", s.Bucket, "/", objectToDelete)

	_, err := s.Client.DeleteObject(ctx, doi)

	if err != nil {
		log.Print("Remove: failed to delete object: s3://", s.Bucket, "/", objectToDelete, err.Error())
	}
	return nil
}

// Stats returns total number of items and their cumulative size
func (s *S3Directory) Stats() (numItems uint64, numBytes uint64) {
	log.Print("Stats: s3 ListObjectsV2 call made for Stats")

	objectCount, sizeOfObjects := GetS3PrefixSize(s.Bucket, s.Prefix)

	return objectCount, sizeOfObjects
}

// GetS3PrefixSize Stats returns total number of items and their cumulative size
func GetS3PrefixSize(bucket, prefix string) (numItems uint64, numBytes uint64) {

	objectCount := uint64(0)
	sizeOfObjects := uint64(0)

	ctx := context.Background()
	params := s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix,
	}

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(endpointResolver))
	if err != nil {
		log.Print("Error loading AWS config: ", err)
	}
	client := s3.NewFromConfig(cfg)

	val, err := client.ListObjectsV2(ctx, &params)
	if err != nil {
		log.Print("Stats: failed to list objects: ", err.Error())
		return 0, 0
	}

	for _, obj := range val.Contents {
		size := uint64(obj.Size)
		objectCount++
		sizeOfObjects += size
	}

	return objectCount, sizeOfObjects
}

// Sync ensures directory metadata itself has been committed
func (s *S3Directory) Sync() error {
	return nil
}

// Lock ensures this process has exclusive access to write in this directory
func (s *S3Directory) Lock() error {
	return nil
}

// Unlock releases the lock held on this directory
func (s *S3Directory) Unlock() error {
	return nil
}
