package main

import (
	"fmt"
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
	if err != nil {
		fmt.Println(err)
	}
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
