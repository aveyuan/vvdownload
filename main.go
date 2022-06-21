package main

import (
	"github.com/go-playground/validator"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
)

func main() {
	InitDB()
	app := iris.New()
	app.Use(recover.New())
	app.Validator = validator.New()
	view := iris.HTML("./template/view", ".html")
	view.Reload(true)
	app.RegisterView(view)
	app.HandleDir("/static", iris.Dir("./template/static"))
	app.HandleDir("/download", iris.Dir("./download"))

	base := NewBase()
	app.Get("/", base.Index)
	app.Post("/downloadFile", base.DownLoadFile)
	app.Delete("/delete", base.DeleteFile)
	app.Get("/list", base.GetFiles)
	app.Listen(":8080")
}
