package main

import (
	"github.com/RaymondCode/simple-demo/mysql"
	"github.com/RaymondCode/simple-demo/pkg/snowflake"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	initRouter(r)
	mysql.InitMysql()
	err := snowflake.Init(1)
	if err != nil {
		log.Fatal("雪花算法初始化失败")
	}
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
