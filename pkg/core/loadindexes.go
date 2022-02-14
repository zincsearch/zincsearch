package core

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"

	"github.com/prabhatsharma/zinc/pkg/zutils"
)

var systemIndexList = []string{"_users", "_index_mapping", "_index_template"}

func LoadZincSystemIndexes() (map[string]*Index, error) {
	log.Print("Loading system indexes...")

	indexList := make(map[string]*Index)
	for _, systemIndex := range systemIndexList {
		tempIndex, err := NewIndex(systemIndex, "disk")
		if err != nil {
			log.Print("Error loading system index: ", systemIndex, " : ", err.Error())
			return nil, err
		}
		indexList[systemIndex] = tempIndex
		indexList[systemIndex].IndexType = "system"
		log.Print("Index loaded: " + systemIndex)
	}

	return indexList, nil
}

func LoadZincIndexesFromDisk() (map[string]*Index, error) {
	log.Print("Loading indexes... from disk")

	indexList := make(map[string]*Index)
	dataPath := zutils.GetEnv("ZINC_DATA_PATH", "./data")
	files, err := os.ReadDir(dataPath)
	if err != nil {
		log.Fatal().Msg("Error reading data directory: " + err.Error())
	}

	for _, f := range files {
		iName := f.Name()
		iNameIsSystemIndex := false
		for _, systemIndex := range systemIndexList {
			if iName == systemIndex {
				iNameIsSystemIndex = true
			}
		}
		if iNameIsSystemIndex {
			continue
		}

		tempIndex, err := NewIndex(iName, "disk")
		if err != nil {
			log.Print("Error loading index: ", iName, " : ", err.Error()) // inform and move in to next index
		} else {
			indexList[iName] = tempIndex
			indexList[iName].IndexType = "user"
			log.Print("Index loaded: " + iName)
		}
	}

	return indexList, nil
}

func LoadZincIndexesFromS3() (map[string]*Index, error) {
	log.Print("Loading indexes from s3...")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Print("Error loading AWS config: ", err)
	}
	client := s3.NewFromConfig(cfg)
	IndexList := make(map[string]*Index)
	dataPath := zutils.GetEnv("ZINC_S3_BUCKET", "")
	delimiter := "/"
	ctx := context.Background()
	params := s3.ListObjectsV2Input{
		Bucket:    &dataPath,
		Delimiter: &delimiter,
	}

	val, err := client.ListObjectsV2(ctx, &params)
	if err != nil {
		log.Print("failed to list indexes in s3: ", err.Error())
		return nil, err
	}

	for _, obj := range val.CommonPrefixes {
		iName := (*obj.Prefix)[0 : len(*obj.Prefix)-1]
		tempIndex, err := NewIndex(iName, "s3")
		if err != nil {
			log.Print("failed to load index "+iName+" in s3: ", err.Error())
		} else {
			IndexList[iName] = tempIndex
			IndexList[iName].IndexType = "user"
			IndexList[iName].StorageType = "s3"
			log.Print("Index loaded: " + iName)
		}

	}

	return IndexList, nil
}

func LoadZincIndexesFromMinIO() (map[string]*Index, error) {
	log.Print("Loading indexes from minio...")

	endpoint := zutils.GetEnv("ZINC_MINIO_ENDPOINT", "")
	accessKeyID := zutils.GetEnv("ZINC_MINIO_ACCESS_KEY_ID", "")
	secretAccessKey := zutils.GetEnv("ZINC_MINIO_SECRET_ACCESS_KEY", "")
	MINIO_BUCKET := zutils.GetEnv("ZINC_MINIO_BUCKET", "")
	if MINIO_BUCKET == "" {
		return nil, nil
	}

	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	}

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, opts)
	if err != nil {
		log.Print(err)
	}

	IndexList := make(map[string]*Index)
	optsList := minio.ListObjectsOptions{
		Recursive: false,
	}
	val := minioClient.ListObjects(context.TODO(), MINIO_BUCKET, optsList)
	if err != nil {
		log.Print("failed to list indexes in minio: ", err.Error())
		return nil, err
	}

	for iName := range val {
		indexName := iName.Key[:len(iName.Key)-1]
		tempIndex, err := NewIndex(indexName, "minio")
		if err != nil {
			log.Print("failed to load index "+iName.Key+" in minio: ", err.Error())
		} else {
			IndexList[indexName] = tempIndex
			IndexList[indexName].IndexType = "user"
			IndexList[indexName].StorageType = "minio"
			log.Print("Index loaded: " + indexName)
		}
	}

	return IndexList, nil
}
