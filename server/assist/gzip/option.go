package gzip

import (
	"net/http"
	"strings"
	"compress/gzip"
	"github.com/gin-gonic/gin"
)

var (
	DefaultExcludeExts = NewExcludedExts([]string{
		".png", ".gif", ".jpeg", ".jpg",
	})
	DefaultOptions = &Options{
		ExcludedExts: DefaultExcludeExts,
	}
)

type Options struct {
	ExcludedExts	ExcludedExts
	ExcludedPaths	ExcludedPaths
	DecompressFn	func(c *gin.Context)
}

type Option func(*Options)

func WithExcludedExts(args []string) Option {
	return func(o *Options) {
		o.ExcludedExts = NewExcludedExts(args)
	}
}

func WithExcludedPaths(args []string) Option {
	return func(o *Options) {
		o.ExcludedPaths = NewExcludedPaths(args)
	}
}

func WithDecompressFn(decompressFn func(c *gin.Context)) Option {
	return func(o *Options) {
		o.DecompressFn = decompressFn
	}
}

type ExcludedExts map[string]bool

// NewExcludedExts 不执行压缩的文件扩展
func NewExcludedExts(exts []string) ExcludedExts {
	res := make(ExcludedExts)
	for _, e := range exts {
		res[e] = true
	}
	return res
}

// Contain 是否包含指定某扩展类型
func (e ExcludedExts) Contains(ext string) bool {
	_, ok := e[ext]
	return ok
}

type ExcludedPaths []string

// NewExcludedPaths 不执行压缩的文件路径
func NewExcludedPaths(paths []string) ExcludedPaths {
	return ExcludedPaths(paths)
}

func (e ExcludedPaths) Contains(path string) bool {
	for _, p := range e {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func DefaultDecompressHandle(c *gin.Context) {
	if c.Request.Body == nil {
		return
	}
	r, err := gzip.NewReader(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.Request.Header.Del("Content-Encoding")
	c.Request.Header.Del("Content-Length")
	c.Request.Body = r
}


