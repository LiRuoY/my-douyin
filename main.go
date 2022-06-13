package main

import (
	"douyin/config"
	"douyin/dao"
	"douyin/pkg/routers"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Init() {
	config.InitConfig("config.yml")
	dao.Init()
}

func main() {
	Init()
	r := gin.Default()
	routers.RegisterRouters(r)
	if err := r.Run(config.CServer.GetPortLikeInDomain()); err != nil {
		fmt.Println(err)
	}
}
