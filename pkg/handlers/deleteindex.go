package handlers

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"

	"github.com/prabhatsharma/zinc/pkg/core"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

// DeleteIndex deletes a zinc index and its associated data. Be careful using thus as you ca't undo this action.
func DeleteIndex(c *gin.Context) {
	indexName := c.Param("target")

	// 0. Check if index exists and Get the index storage type - disk, s3 or memory
	index, exists := core.GetIndex(indexName)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index " + indexName + " does not exists"})
		return
	}

	// 1. Close the index writer
	index.Writer.Close()

	// 2. Delete from the cache
	delete(core.ZINC_INDEX_LIST, index.Name)

	// 3. Physically delete the index
	if index.StorageType == "disk" {
		dataPath := zutils.GetEnv("ZINC_DATA_PATH", "./data")
		err := os.RemoveAll(dataPath + "/" + index.Name)
		if err != nil {
			log.Error().Msgf("failed to delete index: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if index.StorageType == "s3" {
		err := deleteFilesForIndexFromS3(index.Name)
		if err != nil {
			log.Error().Msgf("failed to delete index: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if index.StorageType == "minio" {
		err := deleteFilesForIndexFromMinIO(index.Name)
		if err != nil {
			log.Error().Msgf("failed to delete index: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// delete meta
	if err := core.DeleteIndex(index.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted",
		"index":   index.Name,
		"storage": index.StorageType,
	})
}

func deleteFilesForIndexFromMinIO(indexName string) error {
	endpoint := zutils.GetEnv("ZINC_MINIO_ENDPOINT", "")
	accessKeyID := zutils.GetEnv("ZINC_MINIO_ACCESS_KEY_ID", "")
	secretAccessKey := zutils.GetEnv("ZINC_MINIO_SECRET_ACCESS_KEY", "")
	minioBucket := zutils.GetEnv("ZINC_MINIO_BUCKET", "")

	opts := minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	}

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &opts)
	if err != nil {
		log.Print(err)
	}

	listOpts := minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    indexName,
	}

	objects := minioClient.ListObjects(context.TODO(), minioBucket, listOpts)

	for object := range objects {
		log.Print("Deleting: ", object.Key)

		if err != nil {
			return err
		}
	}

	return nil
}

func deleteFilesForIndexFromS3(indexName string) error {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Print("Error loading AWS config: ", err)
		return err
	}
	client := s3.NewFromConfig(cfg)

	s3bucket := zutils.GetEnv("ZINC_S3_BUCKET", "")
	ctx := context.Background()

	// List Objects in the bucket at prefix
	listObjectsInput := &s3.ListObjectsV2Input{
		Bucket: &s3bucket,
		Prefix: &indexName,
	}
	listObjectsOutput, err := client.ListObjectsV2(ctx, listObjectsInput)
	if err != nil {
		log.Print("failed to list objects: ", err.Error())
		return err
	}

	var fileList []types.ObjectIdentifier

	for _, object := range listObjectsOutput.Contents {
		fileList = append(fileList, types.ObjectIdentifier{
			Key: object.Key,
		})
		log.Print("Deleting: ", *object.Key)
	}

	doi := &s3.DeleteObjectsInput{
		Bucket: &s3bucket,
		Delete: &types.Delete{
			Objects: fileList,
		},
	}
	_, err = client.DeleteObjects(ctx, doi)
	if err != nil {
		log.Print("failed to delete index: ", err.Error())
		return err
	}

	return nil
}
