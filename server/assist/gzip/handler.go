package gzip

import (
	"compress/gzip"
	"io/ioutil"
	"sync"

	"github.com/gin-gonic/gin"
)

type gzipHandler struct {
	*Options
	gzPool sync.Pool
}

func NewGzipHandler(level int, options ...Option) *gzipHandler {
	var gzPool sync.Pool
	gzPool.New = func() interface{} {
		gz, err := gzip.NewWriterLevel(ioutil.Discard, level)
		if err != nil {
			panic(err)
		}
		return gz
	}
	handler := &gzipHandler{
		Options: DefaultOptions,
		gzPool: gzPool,
	}

	for _, setter := range options {
		setter(handler.Options)
	}
	return handler
}

func (g *gzipHandler) Handle(c *gin.Context) {

}