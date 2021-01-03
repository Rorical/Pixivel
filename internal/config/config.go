package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type SQLSetting struct {
	URI  string
	Type string
}

type HashDBSetting struct {
	URI string
}

type RedisSetting struct {
	Host string
	Port string
}

type Setting struct {
	SQL    SQLSetting
	HashDB HashDBSetting
	Redis  RedisSetting
}

func Read() *Setting {
	var settings Setting
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.SetConfigType("json")
	err := v.ReadInConfig()
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	if err := v.Unmarshal(&settings); err != nil {
		fmt.Println(err)
	}
	return &settings
}
