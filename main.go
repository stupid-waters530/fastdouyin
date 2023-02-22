package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ikuraoo/fastdouyin/configure"
	"github.com/ikuraoo/fastdouyin/entity"
	"github.com/ikuraoo/fastdouyin/util"
)

//func main() {
//	util.CreateTable()
//}

func main() {
	configure.InitConfig()
	entity.Init()

	//go service.RunMessageServer()
	util.InitLogger()

	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
