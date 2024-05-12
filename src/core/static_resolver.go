package core

import (
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

var (
	rootHtml   []byte
	staticPath string
)

func init() {
	dirPath, err := os.Getwd()
	if err != nil {
		return
	}

	staticPath = path.Join(dirPath, "static")
	fileBytes, err := os.ReadFile(path.Join(staticPath, "index.html"))
	if err != nil {
		return
	}

	rootHtml = fileBytes
}

func RegisterStaticResolver(router *gin.Engine) {
	router.Static("/assets", path.Join(staticPath, "assets"))

	router.NoRoute(func(ctx *gin.Context) {
		ctx.Data(HttpStatusOK, HTML_MIME_TYPE, rootHtml)
	})
}
