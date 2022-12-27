// Copyright 2022 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package pprof

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
)

func Test_getPrefix(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"default value", nil, "/debug/pprof"},
		{"test user input value", []string{"test/pprof"}, "test/pprof"},
		{"test user input value", []string{"test/pprof", "pprof"}, "test/pprof"},
	}
	for _, tt := range tests {
		if got := getPrefix(tt.args...); got != tt.want {
			t.Errorf("%q. getPrefix() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

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

func Test_Pprof_Router_Group(t *testing.T) {
	bearerToken := "Bearer token"
	h := server.New()
	Register(h)
	adminGroup := h.Group("/admin", func(c context.Context, ctx *app.RequestContext) {
		if ctx.Request.Header.Get("Authorization") != bearerToken {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		ctx.Next(c)
	})
	RouteRegister(adminGroup, "pprof")

	resp := ut.PerformRequest(h.Engine, http.MethodGet, "/admin/pprof/", nil)
	assert.DeepEqual(t, http.StatusForbidden, resp.Code)

	header := ut.Header{
		Key:   "Authorization",
		Value: bearerToken,
	}
	resp = ut.PerformRequest(h.Engine, http.MethodGet, "/admin/pprof/", nil, header)
	assert.DeepEqual(t, http.StatusOK, resp.Code)
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
			assert.DeepEqual(t, http.StatusOK, resp.Code)
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
	assert.DeepEqual(t, http.StatusNotFound, resp.Code)
}
