package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
)

var port string

func init() {
	gin.SetMode(gin.ReleaseMode)
	flag.StringVar(&port, "p", "8080", "-p port")
	flag.Parse()
}

func main() {
	InitDB()
	r := gin.Default()
	r.StaticFS("/download", gin.Dir("./download", false))
	r.LoadHTMLGlob("views/*")
	base := NewBase()
	r.GET("/", base.Index)
	r.GET("/downloadFile", base.DownLoadFile)
	r.GET("/delete", base.DeleteFile)
	r.GET("/list", base.GetFiles)

	log.Print("程序启动成功，监听端口:", port)
	r.Run(":" + port)
}
