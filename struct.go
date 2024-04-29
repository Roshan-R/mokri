package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Item struct {
	Method   string   `yaml:"method"`
	Path     string   `yaml:"path"`
	Response Response `yaml:"response"`
}

type Response struct {
	Status int    `yaml:"status"`
	Body   string `yaml:"body"`
}

type Config struct {
	Routes map[string]Item
}

const configFileName = "config.yaml"
const configDirName = "gomock"

func GetConfigPath() string {
	homeDir, err := os.UserConfigDir()
	Check(err)
	configPath := filepath.Join(homeDir, configDirName, configFileName)
	return configPath
}

func ReadConfig() (Config, error) {
	var config Config
	cfp := GetConfigPath()
	f, err := os.ReadFile(cfp)
	yaml.Unmarshal(f, &config)
	Check(err)

	if config.Routes == nil {
		config.Routes = make(map[string]Item)
	}

	return config, nil
}

func WriteConfig(newConfig *Config) error {
	cfp := GetConfigPath()
	b, err := yaml.Marshal(&newConfig)
	Check(err)

	f, err := os.OpenFile(cfp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	f.Write(b)

	return nil

}
