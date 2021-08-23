package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

var (
	config Config
)

//Config struct
type Config struct {
	AgentNamespace string `json:"agent_namespace"`
	ReconcileTimer int    `json:"reconciletimer"`
	// TODO add more configs
}

// InitConfig loads the config json file
func InitConfig(cfgPath string) error {

	cfgData, err := ioutil.ReadFile(filepath.Clean(cfgPath))
	if err != nil {
		fmt.Printf("Unable read config file" + err.Error())
		return err
	}
	err = json.Unmarshal(cfgData, &config)
	if err != nil {
		fmt.Printf("Failed to load configuration: %v", err)
	}
	return nil
}

// GetConfig returns the config
func GetConfig() *Config {
	return &config
}
