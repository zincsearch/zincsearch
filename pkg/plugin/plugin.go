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
	log.Info().Msg("Loading plugin...")

	dataPath := zutils.GetEnv("ZINC_PLUGIN_PATH", "./plugins")
	files, err := os.ReadDir(dataPath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Error().Msgf("Loading plugin: reading plugin directory err: " + err.Error())
		}
		return
	}

	for _, f := range files {
		log.Info().Msgf("Loading plugin: [%s]", f.Name())
		if !f.IsDir() {
			log.Error().Msgf("Loading plugin: [%s] err: plugin should be a directory", f.Name())
			continue
		}
		p, err := plugin.Open(dataPath + "/" + f.Name() + "/" + f.Name() + ".so")
		if err != nil {
			log.Fatal().Msgf("Loading plugin: [%s] err: %v", dataPath+"/"+f.Name(), err)
			return
		}
		v, err := p.Lookup("ZINC_PLUGIN_TYPE")
		if err != nil {
			log.Fatal().Msgf("Loading plugin: [%s] err: [ZINC_PLUGIN_TYPE] not defined", dataPath+"/"+f.Name())
			return
		}
		pluginType, _ := v.(*string)
		switch *pluginType {
		case TYPE_ANALYSIS:
			err = analysis.Load(p)
		case TYPE_FEATURE:
			err = nil
		default:
			log.Fatal().Msgf("Loading plugin: [%s] err: [ZINC_PLUGIN_TYPE:%s] not supported", dataPath+"/"+f.Name(), *pluginType)
			return
		}
		if err != nil {
			log.Fatal().Msgf("Loading plugin: [%s] err: [ZINC_PLUGIN_TYPE:%s] %v", dataPath+"/"+f.Name(), *pluginType, err)
			return
		}
	}
}
