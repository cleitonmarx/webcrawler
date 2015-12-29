package infrastructure

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/cleitonmarx/webcrawler/config"
)

//ConfigFileRepository implements config.Repository interface
type ConfigFileRepository struct {
	filePath string
}

//GetSystemConfiguration returns the system configurations
func (cfr ConfigFileRepository) GetSystemConfiguration() (config.Configuration, error) {
	conf := config.Configuration{}
	file, e := ioutil.ReadFile(cfr.filePath)
	if e != nil {
		return conf, e
	}

	json.Unmarshal(file, &conf)
	cfr.setEnvConfig(&conf)

	return conf, nil

}

func (cfr ConfigFileRepository) setEnvConfig(currentConfig *config.Configuration) {
	env := os.Getenv("APP_ENV")
	if len(env) > 0 {
		currentConfig.CurrentEnvironment = env
	}
}

//NewConfigFileRepository creates a new config.Repository instance
func NewConfigFileRepository(path string) config.Repository {
	return &ConfigFileRepository{path}
}
