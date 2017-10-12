package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/ieee0824/getenv"
	"os"
	"io/ioutil"
)

var (
	buffer []byte
)

func open(ctx *gin.Context) {
}

func wq(ctx *gin.Context) {
	if err := ioutil.WriteFile(os.Args[1], buffer, 0644); err != nil {
		panic(err)
	}
	fmt.Fprint(ctx.Writer, string(buffer))
	os.Exit(0)
}

func root(ctx *gin.Context) {
	fmt.Fprint(ctx.Writer, string(buffer))
}

func edit(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		panic(err)
	}
	buffer = body
}

func main() {
	var err error
	if len(os.Args) != 2 {
		panic(fmt.Errorf("file name not set"))
	}

	buffer, err = ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.GET("/", root)
	router.POST("/", root)
	router.POST("/o", open)
	router.POST("/wq", wq)
	router.POST("/e", edit)

	if err := router.Run(fmt.Sprintf(":%d", getenv.Int("EDITOR_PORT", 8080))); err != nil {
		panic(err)
	}
}