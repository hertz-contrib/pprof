package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"hertz-contrib-beiye/pprof"
)

func main() {
	h := server.Default()

	pprof.Register(h)

	adminGroup := h.Group("/admin")

	adminGroup.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
	})

	pprof.RouteRegister(adminGroup, "pprof")

	h.Spin()
}
