package web

import "github.com/gin-gonic/gin"

type InitOptions struct {
	Debug bool
}

func Init(opt *InitOptions) {
	if opt.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}
