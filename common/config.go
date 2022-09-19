package common

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func GetConfig(path string) {
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config/")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("Env local not found")
		os.Exit(1)
	}

	viper.SetConfigName(path)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	err = viper.MergeInConfig()
	if err != nil {
		log.Fatalln("Env " + path + " not found")
		os.Exit(1)
	}
}
