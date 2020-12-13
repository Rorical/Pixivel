package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type DatabaseSetting struct {
	URI  string
	Type string
}

var databaseConf = DatabaseSetting{}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Conf error %s\n", err)
		os.Exit(1)
	}
	databaseConf = DatabaseSetting{
		Type: viper.GetString("sql.type"),
		URI:  viper.GetString("sql.uri"),
	}
}

func DatabaseConf() DatabaseSetting {
	return databaseConf
}
