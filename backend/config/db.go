package config

import (
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
	cfgType := cfgVal.Elem()

	for i := 0; i < cfgVal.Type().NumField(); i++ {
		t := cfgType.Type().Field(i).Tag.Get("env")
		if !isEmpty(os.Getenv(t)) {
			if cfgType.Field(i).CanSet() {
				cfgType.Field(i).SetString(os.Getenv(t))
			}
		}
	}
}

func isEmpty(val string) bool {
	return val == ""
}
