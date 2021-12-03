package core

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/rs/zerolog/log"
)

var systemIndexList = []string{"_users", "_index_mapping"}

// var systemIndexList = []string{}

func LoadZincSystemIndexes() (map[string]*Index, error) {
	godotenv.Load()
	// zerolog.TimeFieldFormat = ""
	log.Print("Loading system indexes...")
	// log.Log().Msg("Hello")

	IndexList := make(map[string]*Index)
	var err error

	for _, systemIndex := range systemIndexList {
		IndexList[systemIndex], err = NewIndex(systemIndex)
		IndexList[systemIndex].IndexType = "system"
		if err != nil {
			log.Print(err.Error())
			return nil, err
		}
		log.Print("Index loaded: " + systemIndex)
	}

	return IndexList, nil
}

func LoadZincIndexes() (map[string]*Index, error) {
	godotenv.Load()
	// zerolog.TimeFieldFormat = ""
	log.Print("Loading indexes...")
	// log.Log().Msg("Hello")

	IndexList := make(map[string]*Index)

	DATA_PATH := ""
	if os.Getenv("DATA_PATH") == "" {
		DATA_PATH = "./data"
		// DATA_PATH = "/Users/prabhat/projects/prabhatsharma/zinc/data"
	} else {
		DATA_PATH = os.Getenv("DATA_PATH")
	}

	files, err := os.ReadDir(DATA_PATH)
	if err != nil {
		log.Print("Error reading data directory: ", err.Error())
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

		if !iNameIsSystemIndex {
			tempIndex, err := NewIndex(iName)
			if err != nil {
				log.Print("Error loading index: ", iName, " : ", err.Error()) // inform and move in to next index
			} else {
				IndexList[iName] = tempIndex
				IndexList[iName].IndexType = "user"
				log.Print("Index loaded: " + iName)
			}
		}
	}

	return IndexList, nil
}
