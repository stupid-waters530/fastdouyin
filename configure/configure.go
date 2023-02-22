package configure

import "os"

import (
	"github.com/spf13/viper"
)

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("configure")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/configure")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
