package main

import (
	"fmt"
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

func GetConfigPath() string {
	configDir, err := os.UserConfigDir()
	Check(err)
	configPath := filepath.Join(configDir, "gomock", "config.yaml")
	return configPath
}

func ReadConfig() (Config, error) {
	var config Config
	cfp := GetConfigPath()
	f, err := os.ReadFile(cfp)
	yaml.Unmarshal(f, &config)
	Check(err)
	fmt.Println(config)

	if config.Routes == nil {
		config.Routes = make(map[string]Item)
	}

	return config, nil
}

func WriteConfig(newConfig *Config) error {
	cfp := GetConfigPath()
	y, err := yaml.Marshal(&newConfig)
	Check(err)

	f, err := os.OpenFile(cfp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	f.WriteString(string(y))

	return nil

}
