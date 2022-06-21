package main

import (
	"flag"

	"github.com/go-playground/validator"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8080", "-p port")
	flag.Parse()
}

func main() {
	InitDB()
	app := iris.New()
	app.Use(recover.New())
	app.Validator = validator.New()
	view := iris.HTML("./views", ".html")
	view.Reload(true)
	app.RegisterView(view)
	app.HandleDir("/download", iris.Dir("./download"))

	base := NewBase()
	app.Get("/", base.Index)
	app.Get("/downloadFile", base.DownLoadFile)
	app.Get("/delete", base.DeleteFile)
	app.Get("/list", base.GetFiles)
	app.Listen(":" + port)
}
