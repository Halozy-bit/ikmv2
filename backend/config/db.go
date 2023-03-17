package config

import (
	"log"
	"os"
	"reflect"
)

type MongoConfig struct {
	Driver   string `env:"mongo_driver"`
	User     string `env:"mongo_user"`
	Password string `env:"mongo_password"`
	Address  string `env:"mongo_host"`
	DbName   string `env:"mongo_database"`
}

func AutoEnv(cfg *MongoConfig) {
	cfgVal := reflect.ValueOf(cfg)
	cfgEl := cfgVal.Elem()
	for i := 0; i < cfgEl.Type().NumField(); i++ {
		t := cfgEl.Type().Field(i).Tag.Get("env")
		if os.Getenv(t) != "" {
			log.Printf("get env from %s", t)
			if cfgEl.Field(i).CanSet() {
				cfgEl.Field(i).SetString(os.Getenv(t))
			}
		}
	}
}
