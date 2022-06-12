package main

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/mysql"
	"github.com/RaymondCode/simple-demo/pkg/snowflake"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var (
	fileSystem http.Handler
)

//func init() {
//	// fileSystem = http.FileServer(http.Dir("/var/jenkins_home/workspace/alumni-association/serve/resources/"))
//	fileSystem = http.FileServer(http.Dir("E:\\系统默认\\桌面\\tikTok\\simple-demo\\public"))
//
//}
func main() {
	r := gin.Default()
	initRouter(r)
	mysql.InitMysql()
	err := snowflake.Init(1)
	if err != nil {
		log.Fatal("雪花算法初始化失败")
	}
	// 静态资源请求
	//http.Handle("/public/", http.StripPrefix("/public/", fileSystem))
	if err != nil {
		fmt.Println(err)
	}
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
