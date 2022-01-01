package handlers

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/blugelabs/bluge"
	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/core"
	"github.com/prabhatsharma/zinc/pkg/zutils"
	"github.com/rs/zerolog/log"
)

// DeleteIndex deletes a zinc index and its associated data. Be careful using thus as you ca't undo this action.
func DeleteIndex(c *gin.Context) {
	indexName := c.Param("indexName")

	// 0. Get the index storage type - disk, s3 or memory
	indexStorageType := core.ZINC_INDEX_LIST[indexName].StorageType
	// 1. Close the index writer
	core.ZINC_INDEX_LIST[indexName].Writer.Close()

	// 2. Delete from the cache
	delete(core.ZINC_INDEX_LIST, indexName)

	// 3. Physically delete the index
	deleteIndexMapping := false
	if indexStorageType == "disk" {
		DATA_PATH := zutils.GetEnv("DATA_PATH", "./data")

		err := os.RemoveAll(DATA_PATH + "/" + indexName)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		} else {
			deleteIndexMapping = true
		}
	} else if indexStorageType == "s3" {
		// Load the Shared AWS Configuration (~/.aws/config)
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			log.Print("Error loading AWS config: ", err)
		}
		client := s3.NewFromConfig(cfg)

		S3_BUCKET := zutils.GetEnv("S3_BUCKET", "zinc1")

		ctx := context.Background()
		doi := &s3.DeleteObjectInput{
			Bucket: &S3_BUCKET,
			Key:    &indexName,
		}
		_, err = client.DeleteObject(ctx, doi)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			log.Print("failed to delete index: ", err.Error())
		} else {
			deleteIndexMapping = true
		}
	}

	if deleteIndexMapping {
		// 4. Delete the index mapping
		bdoc := bluge.NewDocument(indexName)
		err := core.ZINC_SYSTEM_INDEX_LIST["_index_mapping"].Writer.Delete(bdoc.ID())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "Deleted",
				"index":   indexName,
			})
		}

	}

}
