package pprof

import (
	"bytes"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_Non_Pprof_Path(t *testing.T) {
	h := server.Default()

	Register(h)

	h.GET("/", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(http.StatusOK, "escaped")
	})

	resp := ut.PerformRequest(h.Engine, http.MethodGet, "/", nil)
	assert.DeepEqual(t, http.StatusOK, resp.Code)

	b, err := ioutil.ReadAll(resp.Body)
	assert.DeepEqual(t, nil, err)
	assert.DeepEqual(t, "escaped", string(b))
}

func Test_Pprof_Index(t *testing.T) {
	h := server.Default()

	Register(h)

	h.GET("/", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(http.StatusOK, "escaped")
	})

	resp := ut.PerformRequest(h.Engine, http.MethodGet, "/debug/pprof/", nil)
	assert.DeepEqual(t, http.StatusOK, resp.Code)
	assert.DeepEqual(t, []byte("text/html; charset=utf-8"), resp.Header().ContentType())

	b, err := ioutil.ReadAll(resp.Body)
	assert.DeepEqual(t, nil, err)
	assert.DeepEqual(t, true, bytes.Contains(b, []byte("<title>/debug/pprof/</title>")))
}

func Test_Pprof_Subs(t *testing.T) {
	h := server.Default()

	Register(h)

	h.GET("/", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(http.StatusOK, "escaped")
	})

	subs := []string{
		"cmdline", "profile", "symbol", "trace", "allocs", "block",
		"goroutine", "heap", "mutex", "threadcreate",
	}

	for _, sub := range subs {
		t.Run(sub, func(t *testing.T) {
			target := "/debug/pprof/" + sub
			if sub == "profile" {
				target += "?seconds=1"
			}
			resp := ut.PerformRequest(h.Engine, http.MethodGet, target, nil)
			assert.DeepEqual(t, 200, resp.Code)
		})
	}
}

func Test_Pprof_Other(t *testing.T) {
	h := server.Default()

	Register(h)

	h.GET("/", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(http.StatusOK, "escaped")
	})

	resp := ut.PerformRequest(h.Engine, http.MethodGet, "/debug/pprof/302", nil)
	assert.DeepEqual(t, 404, resp.Code)
}
