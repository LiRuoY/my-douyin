package config

import (
	"github.com/spf13/viper"
)

var CDB *DB
var CStatic *Static
var CServer *Server

func InitConfig(file string) {
	viper.SetConfigType("yml")
	viper.SetConfigFile(file)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	var myConfig *Config
	if err := viper.Unmarshal(&myConfig); err != nil {
		panic(err)
	}
	CDB = myConfig.DB
	CStatic = myConfig.Static
	CServer = myConfig.Server
}
