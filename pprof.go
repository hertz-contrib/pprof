package pprof

import (
	"net/http/pprof"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/hertz-contrib/pprof/adaptor"
)

const (
	// DefaultPrefix url prefix of pprof
	DefaultPrefix = "/debug/pprof"
)

func getPrefix(prefixOptions ...string) string {
	prefix := DefaultPrefix
	if len(prefixOptions) > 0 {
		prefix = prefixOptions[0]
	}
	return prefix
}

// Register the standard HandlerFuncs from the net/http/pprof package with
// the provided hertz.Hertz. prefixOptions is a optional. If not prefixOptions,
// the default path prefix is used, otherwise first prefixOptions will be path prefix.
func Register(r *server.Hertz, prefixOptions ...string) {
	RouteRegister(&(r.RouterGroup), prefixOptions...)
}

// RouteRegister the standard HandlerFuncs from the net/http/pprof package with
// the provided hertz.RouterGroup. prefixOptions is a optional. If not prefixOptions,
// the default path prefix is used, otherwise first prefixOptions will be path prefix.
func RouteRegister(rg *route.RouterGroup, prefixOptions ...string) {
	prefix := getPrefix(prefixOptions...)

	prefixRouter := rg.Group(prefix)
	{
		prefixRouter.GET("/", adaptor.NewHertzHTTPHandlerFunc(pprof.Index))
		prefixRouter.GET("/cmdline", adaptor.NewHertzHTTPHandlerFunc(pprof.Cmdline))

		prefixRouter.GET("/profile", adaptor.NewHertzHTTPHandlerFunc(pprof.Profile))
		prefixRouter.POST("/symbol", adaptor.NewHertzHTTPHandlerFunc(pprof.Symbol))
		prefixRouter.GET("/symbol", adaptor.NewHertzHTTPHandlerFunc(pprof.Symbol))
		prefixRouter.GET("/trace", adaptor.NewHertzHTTPHandlerFunc(pprof.Trace))
		prefixRouter.GET("/allocs", adaptor.NewHertzHTTPHandlerFunc(pprof.Handler("allocs").ServeHTTP))
		prefixRouter.GET("/block", adaptor.NewHertzHTTPHandlerFunc(pprof.Handler("block").ServeHTTP))
		prefixRouter.GET("/goroutine", adaptor.NewHertzHTTPHandlerFunc(pprof.Handler("goroutine").ServeHTTP))
		prefixRouter.GET("/heap", adaptor.NewHertzHTTPHandlerFunc(pprof.Handler("heap").ServeHTTP))
		prefixRouter.GET("/mutex", adaptor.NewHertzHTTPHandlerFunc(pprof.Handler("mutex").ServeHTTP))
		prefixRouter.GET("/threadcreate", adaptor.NewHertzHTTPHandlerFunc(pprof.Handler("threadcreate").ServeHTTP))
	}
}
