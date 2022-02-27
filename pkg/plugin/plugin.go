package plugin

import (
	"os"
	"plugin"

	"github.com/rs/zerolog/log"

	"github.com/prabhatsharma/zinc/pkg/plugin/analysis"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

const (
	TYPE_ANALYSIS = "analysis"
	TYPE_FEATURE  = "other_feature"
)

func Load() {
	log.Print("Loading plugin... from disk")

	dataPath := zutils.GetEnv("ZINC_PLUGIN_PATH", "./plugin")
	files, err := os.ReadDir(dataPath)
	if err != nil {
		log.Error().Msgf("Error reading plugin directory: " + err.Error())
		return
	}

	for _, f := range files {
		p, err := plugin.Open(dataPath + "/" + f.Name())
		if err != nil {
			log.Fatal().Msgf("Error loading plugin: %s, err %v", dataPath+"/"+f.Name(), err)
			return
		}
		v, err := p.Lookup("ZINC_PLUGIN_TYPE")
		if err != nil {
			log.Fatal().Msgf("Error loading plugin: %s, [ZINC_PLUGIN_TYPE] not defined", dataPath+"/"+f.Name())
			return
		}
		pluginType, _ := v.(*string)
		switch *pluginType {
		case TYPE_ANALYSIS:
			err = analysis.Load(p)
		case TYPE_FEATURE:
			err = nil
		default:
			log.Fatal().Msgf("Error loading plugin: %s, [ZINC_PLUGIN_TYPE:%s] not supported", dataPath+"/"+f.Name(), *pluginType)
			return
		}
		if err != nil {
			log.Fatal().Msgf("Error loading plugin: %s, [ZINC_PLUGIN_TYPE:%s] err %v", dataPath+"/"+f.Name(), *pluginType, err)
			return
		}
	}
}
