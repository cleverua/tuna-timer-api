package config

import (
	"github.com/olebedev/config"
)

func SetupEnvironment(environment string) {

	cfg, err := config.ParseYamlFile("../config.yml")

	if(err != nil) {
		panic("Failed to read/parse config.yml!")
	}

	cfg, err = cfg.Get(environment)

	//host, err := cfg.String("database.host")
	//cfg.Env()
}
