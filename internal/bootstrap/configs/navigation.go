package configs

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type NavConfig struct {
	Layout        string   `yaml:"layout"`
	NavPath       string   `yaml:"nav_path"`
	ErrorFilePath string   `yaml:"error_file_path"`
	FilePath      string   `yaml:"file_path"`
	Pages         []string `yaml:"pages"`
}

func NewNavConfig() *NavConfig {
	cfg := &NavConfig{}
	data, err := os.ReadFile("configs/navbar.yaml")
	if err != nil {
		fmt.Println("Error read navigation file!")
		os.Exit(0)
	}

	if err = yaml.Unmarshal(data, cfg); err != nil {
		fmt.Println("Error read navigation file!")
		os.Exit(0)
	}

	return cfg
}
