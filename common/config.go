package common

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

var ConfigName string

func GetConfig(configName string, fullPath ...string) {
	actualFullPath := "./config/"
	if len(fullPath) != 0 {
		actualFullPath = fullPath[0]
	}
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(actualFullPath)
	viper.AutomaticEnv()

	ConfigName = configName

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("Env local not found")
		os.Exit(1)
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(actualFullPath)
	viper.AutomaticEnv()

	err = viper.MergeInConfig()
	if err != nil {
		log.Fatalln("Env " + configName + " not found")
		os.Exit(1)
	}
}
