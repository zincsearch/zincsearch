package objstore

import (
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"gopkg.in/ini.v1"
	"io"
	"strconv"
)

var ErrBucketNotExists = errors.New("bucket not exists")

type oSSBackend struct {
	client     *oss.Client
	bucketName string
	bucket     *oss.Bucket
}

func createOSSBackend(section *ini.Section) (*oSSBackend, error) {
	key, err := section.GetKey("key_id")
	if err != nil {
		return nil, err
	}
	accessKeyID := key.String()
	key, err = section.GetKey("key")
	if err != nil {
		return nil, err
	}
	accessKeySecret := key.String()
	key, err = section.GetKey("bucket")
	if err != nil {
		return nil, err
	}
	bucketName := key.String()
	key, err = section.GetKey("endpoint")
	if err != nil {
		return nil, err
	}
	endPoint := key.String()

	return newOSSBackend(endPoint, accessKeyID, accessKeySecret, bucketName)
}

func newOSSBackend(endPoint, accessKeyID, accessKeySecret, bucketName string) (*oSSBackend, error) {
	client, err := oss.New(endPoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}

	backend := new(oSSBackend)
	backend.client = client
	backend.bucketName = bucketName

	return backend, nil
}

func (b *oSSBackend) read(key string, w io.Writer) error {
	if b.bucket == nil {
		return ErrBucketNotExists
	}
	objReadCloser, err := b.bucket.GetObject(key)
	if err != nil {
		return err
	}
	defer objReadCloser.Close()

	_, err = io.Copy(w, objReadCloser)
	if err != nil {
		return err
	}

	return nil
}

func (b *oSSBackend) write(key string, r io.Reader) error {
	if b.bucket == nil {
		return ErrBucketNotExists
	}

	err := b.bucket.PutObject(key, io.NopCloser(r))
	if err != nil {
		return err
	}

	return nil
}

func (b *oSSBackend) listObjects(prefix string, accurateCtime bool) ([]ObjectInfo, error) {
	if b.bucket == nil {
		return nil, ErrBucketNotExists
	}

	opts := []oss.Option{oss.Prefix(prefix), oss.MaxKeys(10)}

	var info []ObjectInfo
	for i := 0; ; i++ {
		result, err := b.bucket.ListObjectsV2(opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}
		for _, obj := range result.Objects {
			mtime := obj.LastModified.Unix()
			if accurateCtime {
				mtime, err = b.getAccurateCtime(obj.Key)
				if err != nil {
					mtime = obj.LastModified.Unix()
				}
			}
			info = append(info, ObjectInfo{
				Key: obj.Key, Size: obj.Size,
				LastModified: mtime,
			})
		}

		if !result.IsTruncated {
			break
		} else if i == 0 {
			opts = append(opts, oss.ContinuationToken(result.NextContinuationToken))
		} else {
			opts[len(opts)-1] = oss.ContinuationToken(result.NextContinuationToken)
		}
	}

	return info, nil
}

func (b *oSSBackend) remove(key string) error {
	if b.bucket == nil {
		return ErrBucketNotExists
	}

	err := b.bucket.DeleteObject(key)
	return err
}

func (b *oSSBackend) bucketExists() (bool, error) {
	exists, err := b.client.IsBucketExist(b.bucketName)
	if err != nil {
		return exists, err
	}
	if exists {
		bucket, err := b.client.Bucket(b.bucketName)
		if err != nil {
			return true, err
		}
		b.bucket = bucket
		return true, nil
	}
	return false, nil
}

func (b *oSSBackend) createBucket() error {
	err := b.client.CreateBucket(b.bucketName)
	if err != nil {
		return err
	}

	b.bucket, err = b.client.Bucket(b.bucketName)
	if err != nil {
		return err
	}
	return nil
}

func (b *oSSBackend) getAccurateCtime(key string) (int64, error) {
	meta, err := b.bucket.GetObjectDetailedMeta(key)
	if err != nil {
		return -1, err
	}
	_, ok := meta["X-Oss-Meta-Ctime"]
	if !ok {
		return -1, fmt.Errorf("failed to get X-Oss-Meta-Ctime from metadata")
	}
	ctime, err := strconv.ParseInt(meta["X-Oss-Meta-Ctime"][0], 10, 64)
	if err != nil {
		return -1, err
	}
	return ctime, nil
}
