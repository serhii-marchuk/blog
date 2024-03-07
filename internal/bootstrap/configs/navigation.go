package configs

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type NavItem struct {
	Name        string `yaml:"name"`
	DisplayName string `yaml:"display_name"`
	Uri         string `yaml:"uri"`
	ContentFile string `yaml:"content_file"`
	Active      bool   `yaml:"active,omitempty"`
}

type ErrorPage struct {
	Name        string `yaml:"name"`
	ContentFile string `yaml:"content_file"`
}

type NavConfig struct {
	BaseTemplatePath string      `yaml:"base_template_path"`
	NavBar           []NavItem   `yaml:"nav_bar"`
	ErrorPages       []ErrorPage `yaml:"error_pages"`
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

func (nc *NavConfig) GetContentFileByPageName(name string) string {
	for _, item := range nc.NavBar {
		if item.Name == name {
			return item.ContentFile
		}
	}

	return nc.getErrorPageFile("404")
}

func (nc *NavConfig) getErrorPageFile(name string) string {
	for _, page := range nc.ErrorPages {
		if page.Name == name {
			return page.ContentFile
		}
	}

	return ""
}
