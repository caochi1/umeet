package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type config map[string]string

var (
	Configs map[string]config
	MySQL   config
	Redis   config
	IP      config
)

func init() {
	yamlFile, _ := os.ReadFile("configs.yaml")
	yaml.Unmarshal(yamlFile, &Configs)
	MySQL = Configs["mysql"]
	Redis = Configs["redis"]
	IP = Configs["router"]
}
