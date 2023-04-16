package pprof

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
)

func Test_getFgprofPrefix(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"default value", nil, "/debug/fgprof"},
		{"test user input value", []string{"test/fgprof"}, "test/fgprof"},
		{"test user input value", []string{"test/fgprof", "pprof"}, "test/fgprof"},
	}
	for _, tt := range tests {
		if got := getFgprofPrefix(tt.args...); got != tt.want {
			t.Errorf("%q. getFgprofPrefix() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_Non_Fgprof_Path(t *testing.T) {
	h := server.Default()

	FgprofRegister(h)

	h.GET("/", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(http.StatusOK, "escaped")
	})

	resp := ut.PerformRequest(h.Engine, http.MethodGet, "/", nil)
	assert.DeepEqual(t, http.StatusOK, resp.Code)

	b, err := ioutil.ReadAll(resp.Body)
	assert.DeepEqual(t, nil, err)
	assert.DeepEqual(t, "escaped", string(b))
}

func Test_Fgprof_Index(t *testing.T) {
	h := server.Default()

	FgprofRegister(h)

	h.GET("/", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(http.StatusOK, "escaped")
	})

	resp := ut.PerformRequest(h.Engine, http.MethodGet, "/debug/fgprof/", nil)
	assert.DeepEqual(t, http.StatusOK, resp.Code)
	assert.DeepEqual(t, []byte("application/x-gzip"), resp.Header().ContentType())

	_, err := ioutil.ReadAll(resp.Body)
	assert.DeepEqual(t, nil, err)
}

func Test_Fgprof_Router_Group(t *testing.T) {
	bearerToken := "Bearer token"
	h := server.New()
	// FgprofRegister(h)
	adminGroup := h.Group("/admin", func(c context.Context, ctx *app.RequestContext) {
		if ctx.Request.Header.Get("Authorization") != bearerToken {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		ctx.Next(c)
	})
	FgprofRouteRegister(adminGroup, "fgprof")

	resp := ut.PerformRequest(h.Engine, http.MethodGet, "/admin/fgprof/", nil)
	assert.DeepEqual(t, http.StatusForbidden, resp.Code)

	header := ut.Header{
		Key:   "Authorization",
		Value: bearerToken,
	}
	resp = ut.PerformRequest(h.Engine, http.MethodGet, "/admin/fgprof/", nil, header)
	assert.DeepEqual(t, http.StatusOK, resp.Code)
}

func Test_Fgprof_Other(t *testing.T) {
	h := server.Default()

	FgprofRegister(h)

	h.GET("/", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(http.StatusOK, "escaped")
	})

	resp := ut.PerformRequest(h.Engine, http.MethodGet, "/debug/fpprof/302", nil)
	assert.DeepEqual(t, http.StatusNotFound, resp.Code)
}
