package pprof

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/route"
	"net/http"
	"net/http/pprof"
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
		index := http.HandlerFunc(pprof.Index)
		prefixRouter.GET("/", NewHertzHTTPHandler(index))
		cmdLine := http.HandlerFunc(pprof.Cmdline)
		prefixRouter.GET("/cmdline", NewHertzHTTPHandler(cmdLine))

		prefixRouter.GET("/profile", NewHertzHTTPHandlerFunc(pprof.Profile))
		prefixRouter.POST("/symbol", NewHertzHTTPHandlerFunc(pprof.Symbol))
		prefixRouter.GET("/symbol", NewHertzHTTPHandlerFunc(pprof.Symbol))
		prefixRouter.GET("/trace", NewHertzHTTPHandlerFunc(pprof.Trace))
		prefixRouter.GET("/allocs", NewHertzHTTPHandlerFunc(pprof.Handler("allocs").ServeHTTP))
		prefixRouter.GET("/block", NewHertzHTTPHandlerFunc(pprof.Handler("block").ServeHTTP))
		prefixRouter.GET("/goroutine", NewHertzHTTPHandlerFunc(pprof.Handler("goroutine").ServeHTTP))
		prefixRouter.GET("/heap", NewHertzHTTPHandlerFunc(pprof.Handler("heap").ServeHTTP))
		prefixRouter.GET("/mutex", NewHertzHTTPHandlerFunc(pprof.Handler("mutex").ServeHTTP))
		prefixRouter.GET("/threadcreate", NewHertzHTTPHandlerFunc(pprof.Handler("threadcreate").ServeHTTP))
	}
}
