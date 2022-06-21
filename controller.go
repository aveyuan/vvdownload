package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

const DownloadFileDir = "download"

type Base struct {
	mux sync.Mutex
}

func NewBase() *Base {
	return new(Base)
}

//  ReJson 统一返回调度模型
func (t *Base) ReJson(code int, msg string, data interface{}, err interface{}, c iris.Context) {
	c.JSON(map[string]interface{}{
		"code": code,
		"data": data,
		"msg":  msg,
		"err":  err,
	})
}

//  ReData 只对数据返回
func (t *Base) ReData(code int, data interface{}, c iris.Context) {
	t.ReJson(200, "", data, nil, c)
}

// ReError 只对错误返回
func (t *Base) ReError(code int, msg string, error interface{}, c iris.Context) {
	t.ReJson(code, msg, nil, error, c)
}

// ReOk 操作正确的数据
func (t *Base) ReOk(code int, msg string, c iris.Context) {
	t.ReJson(code, msg, nil, nil, c)
}

// Index 首页数据
func (t *Base) Index(ctx iris.Context) {
	ctx.View("index.html")
}

// GetFiles 获取所有数据
func (t *Base) GetFiles(ctx iris.Context) {
	all, err := new(File).GetFileAll()
	if err != nil {
		t.ReError(402, "数据获取失败", err, ctx)
		return
	}
	t.ReData(200, all, ctx)

}

// DeleteFile 删除一个文件
func (t *Base) DeleteFile(ctx iris.Context) {
	var form struct {
		ID uint `json:"id" form:"id" validate:"required"`
	}
	if err := ctx.ReadForm(&form); err != nil || form.ID == 0 {
		t.ReError(402, "参数错误", err, ctx)
		return
	}
	// 查询到路径，然后删除
	one, err := new(File).Get(form.ID)
	if err != nil {
		t.ReError(500, "未找到数据", err, ctx)
		return
	}
	// 判断文件是否存在，不存在删除数据库
	_, err = os.Stat(one.Path)
	if err != nil {
		if err := new(File).Delete(form.ID); err != nil {
			t.ReError(500, "数据删除失败", err, ctx)
			return
		}
		t.ReOk(200, "删除成功", ctx)
		return
	}
	// 正常流程
	if err := os.Remove(one.Path); err != nil {
		t.ReError(500, "文件删除失败", err, ctx)
		return
	}
	if err := new(File).Delete(form.ID); err != nil {
		t.ReError(500, "数据删除失败", err, ctx)
		return
	}
	t.ReOk(200, "删除成功", ctx)

}

// DownLoadFile 下载一个文件
func (t *Base) DownLoadFile(ctx iris.Context) {
	t.mux.Lock()
	defer t.mux.Unlock()
	var form struct {
		URL  string `json:"url" form:"url" validate:"required"`
		Name string `json:"name" form:"name"  validate:"required"`
	}
	if err := ctx.ReadForm(&form); err != nil {
		t.ReError(402, "参数错误", err, ctx)
		return
	}

	path, err := Download(form.URL)
	if err != nil {
		t.ReError(402, "下载失败", err, ctx)
		return
	}
	f, err := os.Stat(path)
	if err != nil {
		t.ReError(402, "数据获取失败", err, ctx)
		return
	}
	// 创建数据，然后提交到任务
	var one File
	one.Srouce = form.URL
	one.Path = path
	one.Size = f.Size()
	one.Name = form.Name
	if err := one.Create(); err != nil {
		t.ReError(500, "数据存储失败", err, ctx)
		return
	}
	t.ReOk(200, path, ctx)

}

func Download(URL string) (string, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	_, err = os.Stat(DownloadFileDir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(DownloadFileDir, 0755); err != nil {
			return "", err
		}
	}

	saveDir := DownloadFileDir + "/" + uuid.NewString() + filepath.Ext(URL)

	f, err := os.OpenFile(saveDir, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return "", err
	}

	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return "", err
	}

	return saveDir, nil
}
