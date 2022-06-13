package dao

import (
	"douyin/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	//constants.MYSQLDefaultDNS
	DB, err = gorm.Open(mysql.Open(config.CDB.GetMySQLDNS()), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}
}
