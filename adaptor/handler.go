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

package adaptor

import (
	"context"
	"net/http"
	"unsafe"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// NewHertzHTTPHandlerFunc wraps net/http handler to hertz app.HandlerFunc,
// so it can be passed to hertz server
//
// While this function may be used for easy switching from net/http to hertz,
// it has the following drawbacks comparing to using manually written hertz
// request handler:
//
//   - A lot of useful functionality provided by hertz is missing
//     from net/http handler.
//
//   - net/http -> hertz handler conversion has some overhead,
//     so the returned handler will be always slower than manually written
//     hertz handler.
//
// So it is advisable using this function only for net/http -> hertz switching.
// Then manually convert net/http handlers to hertz handlers
func NewHertzHTTPHandlerFunc(h http.HandlerFunc) app.HandlerFunc {
	return NewHertzHTTPHandler(h)
}

// NewHertzHTTPHandler wraps net/http handler to hertz app.HandlerFunc,
// so it can be passed to hertz server
//
// While this function may be used for easy switching from net/http to hertz,
// it has the following drawbacks comparing to using manually written hertz
// request handler:
//
//   - A lot of useful functionality provided by hertz is missing
//     from net/http handler.
//
//   - net/http -> hertz handler conversion has some overhead,
//     so the returned handler will be always slower than manually written
//     hertz handler.
//
// So it is advisable using this function only for net/http -> hertz switching.
// Then manually convert net/http handlers to hertz handlers
func NewHertzHTTPHandler(h http.Handler) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		req, err := adaptor.GetCompatRequest(&c.Request)
		if err != nil {
			hlog.CtxErrorf(ctx, "HERTZ: Get request error: %v", err)
			c.String(http.StatusInternalServerError, consts.StatusMessage(http.StatusInternalServerError))
			return
		}
		req.RequestURI = b2s(c.Request.RequestURI())
		rw := adaptor.GetCompatResponseWriter(&c.Response)
		c.ForEachKey(func(k string, v interface{}) {
			ctx = context.WithValue(ctx, k, v)
		})
		h.ServeHTTP(rw, req.WithContext(ctx))
		body := c.Response.Body()
		// From net/http.ResponseWriter.Write:
		// If the Header does not contain a Content-Type line, Write adds a Content-Type set
		// to the result of passing the initial 512 bytes of written data to DetectContentType.
		if len(c.GetHeader(consts.HeaderContentType)) == 0 {
			l := 512
			if len(body) < 512 {
				l = len(body)
			}
			c.Response.Header.Set(consts.HeaderContentType, http.DetectContentType(body[:l]))
		}
	}
}

// b2s converts byte slice to a string without memory allocation.
// See https://groups.google.com/forum/#!msg/Golang-Nuts/ENgbUzYvCuU/90yGx7GUAgAJ .
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
