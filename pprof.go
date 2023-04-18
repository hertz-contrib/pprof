/*
 * Copyright 2022 CloudWeGo Authors
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
