package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Dir        string   `yaml:"dir"`
	Buildflags []string `yaml:"build_flags"`
	RunFlags   []string `yaml:"run_flags"`
	Ignore     []string `yamll:"ignore"`
	Verbose    bool     `yaml:"verbose"`
}

func loadYmlConfig(cfg *Config, ymlFile string) error {
	file, err := ioutil.ReadFile(ymlFile)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(file, cfg); err != nil {
		return err
	}
	return nil
}

func LoadYml(configFile string) (Config, error) {
	var cfg Config
	if err := loadYmlConfig(&cfg, configFile); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
