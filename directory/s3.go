package directory

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	blugeIndex "github.com/blugelabs/bluge/index"
	segment "github.com/blugelabs/bluge_segment_api"
	zerolog "github.com/rs/zerolog/log"
)

type S3Directory struct {
	Bucket string
	Prefix string
	Client *s3.Client
}

func (s *S3Directory) NewS3Directory(bucket, prefix string) S3Directory {

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		zerolog.Print("Error loading AWS config: ", err)
	}
	client := s3.NewFromConfig(cfg)

	directory := &S3Directory{
		Bucket: bucket,
		Prefix: prefix,
		Client: client,
	}

	return *directory
}

func (s *S3Directory) createObjectName(kind string, id uint64) string {
	return fmt.Sprintf("%012x", id) + kind
}

func (s *S3Directory) Setup(readOnly bool) error {
	return nil
}

// List the ids of all the items of the specified kind
// Items are returned in descending order by id
func (s *S3Directory) List(kind string) ([]uint64, error) {
	var itemList []uint64

	ctx := context.Background()
	params := s3.ListObjectsV2Input{
		Bucket: &s.Bucket,
		Prefix: &s.Prefix,
	}

	val, err := s.Client.ListObjectsV2(ctx, &params)
	if err != nil {
		zerolog.Print("failed to list objects: ", err.Error())
		return nil, err
	}

	for _, obj := range val.Contents {
		if filepath.Ext(*obj.Key) != kind {
			continue
		}

		stringID := filepath.Base(*obj.Key)
		stringID = stringID[:len(stringID)-len(kind)-1]

		parsedID, err := strconv.ParseUint(stringID, 16, 64)
		if err != nil {
			zerolog.Print("failed to parse object id: ", err.Error())
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

	key := s.createObjectName(kind, id)

	goi := &s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	}

	output, err := s.Client.GetObject(ctx, goi)

	if err != nil {
		zerolog.Print("failed to get object: ", err.Error())
		return nil, nil, err
	}

	data, err := ioutil.ReadAll(output.Body)
	if err != nil {
		zerolog.Print("failed to read object", err.Error())
		return nil, nil, err
	}

	return segment.NewDataBytes(data), nil, nil
}

// Persist a new item with data from the provided WriterTo
// Implementations should monitor the closeCh and return with error
// in the event it is closed before completion.
func (s *S3Directory) Persist(kind string, id uint64, w blugeIndex.WriterTo, closeCh chan struct{}) error {
	var buf bytes.Buffer
	_, err := w.WriteTo(&buf, closeCh)
	if err != nil {
		zerolog.Print("failed to write object to buffer: ", err.Error())
		return err
	}

	s3ObjectName := s.createObjectName(kind, id)

	path := filepath.Join(s.Prefix, s3ObjectName)
	reader := bufio.NewReader(&buf)

	ctx := context.Background()

	params := s3.PutObjectInput{
		Bucket: &s.Bucket,
		Key:    &path,
		Body:   reader,
	}

	ouput, err := s.Client.PutObject(ctx, &params)

	if err != nil {
		zerolog.Print("failed to write object: ", err.Error())
		return err
	}

	zerolog.Print("s3 object written. Its md5 hash is: ", ouput.ETag) // TODO: compare md5 hashes here to ensure successful write

	return nil
}

// Remove the specified item
func (s *S3Directory) Remove(kind string, id uint64) error {
	objectToDelete := filepath.Join(s.Prefix, s.createObjectName(kind, id))
	ctx := context.Background()
	doi := &s3.DeleteObjectInput{
		Bucket: &s.Bucket,
		Key:    &objectToDelete,
	}
	_, err := s.Client.DeleteObject(ctx, doi)

	if err != nil {
		zerolog.Print("failed to delete object: ", err.Error())
	}
	return nil
}

// Stats returns total number of items and their cumulative size
func (s *S3Directory) Stats() (numItems uint64, numBytes uint64) {
	objectCount := uint64(0)
	sizeOfObjects := uint64(0)

	ctx := context.Background()
	params := s3.ListObjectsV2Input{
		Bucket: &s.Bucket,
		Prefix: &s.Prefix,
	}

	val, err := s.Client.ListObjectsV2(ctx, &params)
	if err != nil {
		zerolog.Print("failed to list objects: ", err.Error())
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
