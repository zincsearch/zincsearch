package core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/blugelabs/bluge"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"

	zincanalysis "github.com/prabhatsharma/zinc/pkg/uquery/v2/analysis"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

var systemIndexList = []string{"_index_mapping", "_index_template", "_index", "_metadata", "_users"}

func LoadZincSystemIndexes() (map[string]*Index, error) {
	indexList := make(map[string]*Index)
	for _, index := range systemIndexList {
		log.Info().Msgf("Loading system index... [%s:%s]", index, "disk")
		writer, err := LoadIndexWriter(index, "disk")
		if err != nil {
			return nil, err
		}
		indexList[index] = &Index{
			Name:        index,
			IndexType:   "system",
			StorageType: "disk",
			Writer:      writer,
		}
	}

	return indexList, nil
}

func LoadZincIndexesFromMeta() (map[string]*Index, error) {
	query := bluge.NewMatchAllQuery()
	searchRequest := bluge.NewAllMatches(query).WithStandardAggregations()
	reader, _ := ZINC_SYSTEM_INDEX_LIST["_index"].Writer.Reader()
	defer reader.Close()

	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		return nil, fmt.Errorf("core.LoadZincIndexesFromMeta: error executing search: %v", err)
	}

	indexList := make(map[string]*Index)
	next, err := dmi.Next()
	for err == nil && next != nil {
		index := &Index{IndexType: "user"}
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			switch field {
			case "name":
				index.Name = string(value)
			case "index_type":
				index.IndexType = string(value)
			case "storage_type":
				index.StorageType = string(value)
			case "settings":
				json.Unmarshal(value, &index.Settings)
			case "mappings":
				json.Unmarshal(value, &index.CachedMappings)
			default:
			}
			return true
		})

		log.Info().Msgf("Loading user   index... [%s:%s]", index.Name, index.StorageType)
		if err != nil {
			log.Printf("core.LoadZincIndexesFromMeta: error accessing stored fields: %v", err)
		}

		// load index analysis
		if index.Settings != nil && index.Settings.Analysis != nil {
			index.CachedAnalyzers, err = zincanalysis.RequestAnalyzer(index.Settings.Analysis)
			if err != nil {
				log.Printf("core.LoadZincIndexesFromMeta: error parse stored analysis: %v", err)
			}
		}

		// load index data
		index.Writer, err = LoadIndexWriter(index.Name, index.StorageType)
		if err != nil {
			log.Error().Msgf("Loading user   index... [%s:%s] error: %v", index.Name, index.StorageType, err)
		}

		indexList[index.Name] = index

		next, err = dmi.Next()
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

		tempIndex, err := NewIndex(iName, "disk", NotCompatibleNewIndexMeta)
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

	dataPath := zutils.GetEnv("ZINC_S3_BUCKET", "")
	if dataPath == "" {
		return nil, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Print("Error loading AWS config: ", err)
	}
	client := s3.NewFromConfig(cfg)
	IndexList := make(map[string]*Index)
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
		tempIndex, err := NewIndex(iName, "s3", NotCompatibleNewIndexMeta)
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
	dataPath := zutils.GetEnv("ZINC_MINIO_BUCKET", "")
	if dataPath == "" {
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

	val := minioClient.ListObjects(context.TODO(), dataPath, optsList)
	if err != nil {
		log.Print("failed to list indexes in minio: ", err.Error())
		return nil, err
	}

	for iName := range val {
		indexName := iName.Key[:len(iName.Key)-1]
		tempIndex, err := NewIndex(indexName, "minio", NotCompatibleNewIndexMeta)
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
