package objstore

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gopkg.in/ini.v1"
	"io"
	"strconv"
)

type s3Backend struct {
	client     *minio.Client
	bucketName string
}

func createS3Backend(section *ini.Section) (*s3Backend, error) {
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
	var endPoint string
	if key, err = section.GetKey("host"); err == nil {
		endPoint = key.String()
	}
	var useV4Sig bool
	if key, err = section.GetKey("use_v4_signature"); err == nil {
		if ok, err := key.Bool(); err == nil {
			useV4Sig = ok
		}
	}
	var useHTTPS bool
	if key, err = section.GetKey("use_https"); err == nil {
		if ok, err := key.Bool(); err == nil {
			useHTTPS = ok
		}
	}
	var pathStyleRequest bool
	if key, err = section.GetKey("path_style_request"); err == nil {
		if ok, err := key.Bool(); err == nil {
			pathStyleRequest = ok
		}
	}
	var awsRegion string
	if useV4Sig {
		key, err = section.GetKey("aws_region")
		if err != nil {
			return nil, err
		}
		awsRegion = key.String()
	}
	return newS3Backend(
		accessKeyID, accessKeySecret, bucketName,
		endPoint, awsRegion, useHTTPS, pathStyleRequest, useV4Sig)
}

func newS3Backend(accessKeyID, accessKeySecret, bucketName, endPoint, region string, useHTTPS, pathStyleRequest, useV4Sig bool) (*s3Backend, error) {
	var bucketLookup minio.BucketLookupType
	if pathStyleRequest {
		bucketLookup = minio.BucketLookupPath
	}
	if endPoint == "" {
		if region != "" {
			endPoint = "s3." + region + ".amazonaws.com"
		} else {
			endPoint = "s3.amazonaws.com"
		}
	}
	if useV4Sig {
		cli, err := minio.New(endPoint, &minio.Options{
			Creds:        credentials.NewStaticV4(accessKeyID, accessKeySecret, ""),
			Secure:       useHTTPS,
			Region:       region,
			BucketLookup: bucketLookup,
		})
		return &s3Backend{cli, bucketName}, err
	}
	cli, err := minio.New(endPoint, &minio.Options{
		Creds:        credentials.NewStaticV2(accessKeyID, accessKeySecret, ""),
		Secure:       useHTTPS,
		BucketLookup: bucketLookup,
	})
	return &s3Backend{cli, bucketName}, err
}

func (b *s3Backend) read(key string, w io.Writer) error {
	obj, err := b.client.GetObject(context.TODO(), b.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer obj.Close()
	_, err = io.Copy(w, obj)
	if err != nil {
		return err
	}
	return nil
}

func (b *s3Backend) write(key string, r io.Reader) error {
	opts := minio.PutObjectOptions{}
	return b.writeWithMeta(key, opts, r)
}

func (b *s3Backend) writeWithMeta(key string, opts minio.PutObjectOptions, r io.Reader) error {

	_, err := b.client.PutObject(context.TODO(), b.bucketName, key, r, -1, opts)
	if err != nil {
		return err
	}
	return nil
}

func (b *s3Backend) listObjects(prefix string, accurateCtime bool) ([]ObjectInfo, error) {
	opts := minio.ListObjectsOptions{Prefix: prefix, Recursive: true}
	objs := b.client.ListObjects(context.Background(), b.bucketName, opts)

	var info []ObjectInfo
	var err error
	for obj := range objs {
		if obj.Err != nil {
			return nil, fmt.Errorf("failed to list folder: %w", obj.Err)
		}
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

	return info, nil
}

func (b *s3Backend) remove(key string) error {
	opts := minio.RemoveObjectOptions{}
	err := b.client.RemoveObject(context.Background(), b.bucketName, key, opts)
	return err
}

func (b *s3Backend) bucketExists() (bool, error) {
	return b.client.BucketExists(context.Background(), b.bucketName)
}

func (b *s3Backend) createBucket() error {
	opts := minio.MakeBucketOptions{}
	return b.client.MakeBucket(context.Background(), b.bucketName, opts)
}

func (b *s3Backend) getAccurateCtime(key string) (int64, error) {
	objInfo, err := b.client.StatObject(context.TODO(), b.bucketName, key, minio.StatObjectOptions{})
	if err != nil {
		return -1, err
	}
	_, ok := objInfo.Metadata["X-Amz-Meta-Ctime"]
	if !ok {
		return -1, fmt.Errorf("failed to get X-Amz-Meta-Ctime from metadata")
	}
	ctime, err := strconv.ParseInt(objInfo.Metadata["X-Amz-Meta-Ctime"][0], 10, 64)
	if err != nil {
		return -1, err
	}
	return ctime, nil
}
