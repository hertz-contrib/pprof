/*
 * Copyright 2023 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * The MIT License (MIT)
 *
 * Copyright (c) 2015-present Aliaksandr Valialkin, VertaMedia, Kirill Danshin, Erik Dubbelboer, FastHTTP Authors
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 * This file may have been modified by CloudWeGo authors. All CloudWeGo
 * Modifications are Copyright 2022 CloudWeGo Authors.
 */

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
