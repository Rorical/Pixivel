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
	URI         string
	MaxIdle     int
	IdleTimeout int
	Password    string
}

type PixivSetting struct {
	RefreshToken string
	AccessToken  string
}

type FilterSetting struct {
	File string
}

type Setting struct {
	SQL    SQLSetting
	HashDB HashDBSetting
	Redis  RedisSetting
	Pixiv  PixivSetting
	Filter FilterSetting
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
}

func Read() *Setting {
	var settings Setting
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	if err := viper.Unmarshal(&settings); err != nil {
		fmt.Println(err)
	}
	return &settings
}
