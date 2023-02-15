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
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/index"
	segment "github.com/blugelabs/bluge_segment_api"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zincsearch/pkg/config"
)

// GetMinIOConfig returns a bluge config that will store index data in MinIO
// bucket: the MinIO bucket to use
// indexName: the name of the index to use. It will be an MinIO prefix (folder)
func GetMinIOConfig(bucket string, indexName string, timeRange ...int64) bluge.Config {
	config := index.DefaultConfigWithDirectory(func() index.Directory {
		return NewMinIODirectory(bucket, indexName)
	})
	config = config.WithPersisterNapTimeMSec(50)
	if len(timeRange) == 2 {
		if timeRange[0] <= timeRange[1] {
			config = config.WithTimeRange(timeRange[0], timeRange[1])
		}
	}
	return bluge.DefaultConfigWithIndexConfig(config)
}

type MinIODirectory struct {
	Bucket string
	Prefix string
	Client minio.Client
}

// NewMinIODirectory creates a new MinIODirectory instance which can be used to create MinIO backed indexes
func NewMinIODirectory(bucket, prefix string) index.Directory {

	endpoint := config.Global.MinIO.Endpoint
	accessKeyID := config.Global.MinIO.AccessKeyID
	secretAccessKey := config.Global.MinIO.SecretAccessKey

	opts := minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	}

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &opts)
	if err != nil {
		log.Print(err)
	}

	directory := &MinIODirectory{
		Bucket: bucket,
		Prefix: prefix,
		Client: *minioClient,
	}

	return directory
}

func (s *MinIODirectory) fileName(kind string, id uint64) string {
	return fmt.Sprintf("%012x", id) + kind
}

func (s *MinIODirectory) Setup(readOnly bool) error {
	return nil
}

// List the ids of all the items of the specified kind
// Items are returned in descending order by id
func (s *MinIODirectory) List(kind string) ([]uint64, error) {
	log.Print("List: MinIO ListObjects call made: MinIO://", s.Bucket+"/"+s.Prefix)
	var itemList []uint64

	opts := minio.ListObjectsOptions{
		Recursive: true,
	}

	val := s.Client.ListObjects(context.TODO(), s.Bucket, opts)

	for obj := range val {
		fileKind := filepath.Ext(obj.Key)
		if fileKind != kind {
			continue
		}

		stringID := filepath.Base(obj.Key)
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
func (s *MinIODirectory) Load(kind string, id uint64) (*segment.Data, io.Closer, error) {

	// key := s.Prefix + "/" + s.fileName(kind, id)
	key := filepath.Join(s.Prefix, s.fileName(kind, id))

	log.Print("Load: MinIO GetObject call made. MinIO://", s.Bucket, "/", key)

	reader, err := s.Client.GetObject(context.TODO(), s.Bucket, key, minio.GetObjectOptions{})

	if err != nil {
		log.Print("Load: failed to get object: MinIO://", s.Bucket, "/", key, err.Error())
	}
	// defer reader.Close()

	if err != nil {
		log.Print("Load: failed to get object: MinIO://"+s.Bucket+"/"+key, err.Error())
		return nil, nil, err
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Print("Load: failed to read object: ", err.Error())
		return nil, nil, err
	}

	return segment.NewDataBytes(data), nil, nil
}

// Persist a new item with data from the provided WriterTo
// Implementations should monitor the closeCh and return with error
// in the event it is closed before completion.
func (s *MinIODirectory) Persist(kind string, id uint64, w index.WriterTo, closeCh chan struct{}) error {
	var buf bytes.Buffer
	size, err := w.WriteTo(&buf, closeCh)
	if err != nil {
		log.Print("Persist: failed to write object to buffer: ", err.Error())
		return err
	}

	reader := bufio.NewReader(&buf)

	key := filepath.Join(s.Prefix, s.fileName(kind, id))

	output, err := s.Client.PutObject(context.TODO(), s.Bucket, key, reader, size, minio.PutObjectOptions{})

	if err != nil {
		log.Print("Persist: failed to write object: ", err.Error())
		return err
	}

	h := md5.New()
	h.Write(buf.Bytes())

	if output.ETag != fmt.Sprintf("%x", h.Sum(nil)) {
		log.Print("Warning: MinIO object " + s.Bucket + "/" + key + " has incorrect checksum")
	}

	log.Print("Persist: MinIO object "+s.Bucket+"/"+key+" written. Its md5 hash is: ", output.ETag)

	return nil
}

// Remove the specified item
func (s *MinIODirectory) Remove(kind string, id uint64) error {
	objectToDelete := filepath.Join(s.Prefix, s.fileName(kind, id))

	log.Print("Remove: MinIO DeleteObject call made MinIO://", s.Bucket, "/", objectToDelete)

	err := s.Client.RemoveObject(context.TODO(), s.Bucket, objectToDelete, minio.RemoveObjectOptions{})

	if err != nil {
		log.Print("Remove: failed to delete object: MinIO://", s.Bucket, "/", objectToDelete, err.Error())
	}
	return nil
}

// Stats returns total number of items and their cumulative size
func (s *MinIODirectory) Stats() (numItems uint64, numBytes uint64) {
	log.Print("Stats: MinIO ListObjectsV2 call made for Stats")

	objectCount := uint64(0)
	sizeOfObjects := uint64(0)

	log.Print("Stats: MinIO ListObjectsV2 call made for Stats MinIO://", s.Bucket+"/"+s.Prefix)

	opts := minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    s.Prefix,
	}

	objects := s.Client.ListObjects(context.TODO(), s.Bucket, opts)

	for obj := range objects {
		size := uint64(obj.Size)
		objectCount++
		sizeOfObjects += size
	}

	return objectCount, sizeOfObjects
}

// Sync ensures directory metadata itself has been committed
func (s *MinIODirectory) Sync() error {
	return nil
}

// Lock ensures this process has exclusive access to write in this directory
func (s *MinIODirectory) Lock() error {
	return nil
}

// Unlock releases the lock held on this directory
func (s *MinIODirectory) Unlock() error {
	return nil
}
