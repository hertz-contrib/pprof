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
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
)

func TestNewHertzHandler(t *testing.T) {
	t.Parallel()

	expectedMethod := consts.MethodPost
	expectedProto := "HTTP/1.1"
	expectedProtoMajor := 1
	expectedProtoMinor := 1
	expectedRequestURI := "http://foobar.com/foo/bar?baz=123"
	expectedBody := "<!doctype html><html>"
	expectedContentLength := len(expectedBody)
	expectedHost := "foobar.com"
	expectedHeader := map[string]string{
		"Foo-Bar":         "baz",
		"Abc":             "defg",
		"XXX-Remote-Addr": "123.43.4543.345",
	}
	expectedURL, err := url.ParseRequestURI(expectedRequestURI)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedContextKey := "contextKey"
	expectedContextValue := "contextValue"
	expectedContentType := "text/html; charset=utf-8"

	callsCount := 0
	nethttpH := func(w http.ResponseWriter, r *http.Request) {
		callsCount++
		assert.Assertf(t, r.Method == expectedMethod, "unexpected method %q. Expecting %q", r.Method, expectedMethod)
		assert.Assertf(t, r.Proto == expectedProto, "unexpected proto %q. Expecting %q", r.Proto, expectedProto)
		assert.Assertf(t, r.ProtoMajor == expectedProtoMajor, "unexpected protoMajor %d. Expecting %d", r.ProtoMajor, expectedProtoMajor)
		assert.Assertf(t, r.ProtoMinor == expectedProtoMinor, "unexpected protoMinor %d. Expecting %d", r.ProtoMinor, expectedProtoMinor)
		assert.Assertf(t, r.RequestURI == expectedRequestURI, "unexpected requestURI %q. Expecting %q", r.RequestURI, expectedRequestURI)
		assert.Assertf(t, r.ContentLength == int64(expectedContentLength), "unexpected contentLength %d. Expecting %d", r.ContentLength, expectedContentLength)
		assert.Assertf(t, len(r.TransferEncoding) == 0, "unexpected transferEncoding %q. Expecting []", r.TransferEncoding)
		assert.Assertf(t, r.Host == expectedHost, "unexpected host %q. Expecting %q", r.Host, expectedHost)
		body, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			t.Fatalf("unexpected error when reading request body: %v", err)
		}
		assert.Assertf(t, string(body) == expectedBody, "unexpected body %q. Expecting %q", body, expectedBody)
		assert.Assertf(t, reflect.DeepEqual(r.URL, expectedURL), "unexpected URL: %#v. Expecting %#v", r.URL, expectedURL)
		assert.Assertf(t, r.Context().Value(expectedContextKey) == expectedContextValue,
			"unexpected context value for key %q. Expecting %q, in fact: %v", expectedContextKey,
			expectedContextValue, r.Context().Value(expectedContextKey))
		for k, expectedV := range expectedHeader {
			v := r.Header.Get(k)
			if v != expectedV {
				t.Fatalf("unexpected header value %q for key %q. Expecting %q", v, k, expectedV)
			}
		}
		w.Header().Set("Header1", "value1")
		w.Header().Set("Header2", "value2")
		w.WriteHeader(http.StatusBadRequest) // nolint:errcheck
		w.Write(body)
	}
	hertzH := NewHertzHTTPHandler(http.HandlerFunc(nethttpH))
	hertzH = setContextValueMiddleware(hertzH, expectedContextKey, expectedContextValue)
	var ctx app.RequestContext
	var req protocol.Request
	req.Header.SetMethod(expectedMethod)
	req.SetRequestURI(expectedRequestURI)
	req.Header.SetHost(expectedHost)
	req.BodyWriter().Write([]byte(expectedBody)) // nolint:errcheck
	for k, v := range expectedHeader {
		req.Header.Set(k, v)
	}
	req.CopyTo(&ctx.Request)
	hertzH(context.Background(), &ctx)
	assert.Assertf(t, callsCount == 1, "unexpected callsCount: %d. Expecting 1", callsCount)
	resp := &ctx.Response
	assert.Assertf(t, resp.StatusCode() == http.StatusBadRequest, "unexpected statusCode: %d. Expecting %d", resp.StatusCode(), http.StatusBadRequest)
	assert.Assertf(t, string(resp.Header.Peek("Header1")) == "value1", "unexpected header value: %q. Expecting %q", resp.Header.Peek("Header1"), "value1")
	assert.Assertf(t, string(resp.Header.Peek("Header2")) == "value2", "unexpected header value: %q. Expecting %q", resp.Header.Peek("Header2"), "value2")
	assert.Assertf(t, string(resp.Body()) == expectedBody, "unexpected response body %q. Expecting %q", resp.Body(), expectedBody)
	assert.Assertf(t, string(resp.Header.Peek("Content-Type")) == expectedContentType, "unexpected content-type %q. Expecting %q", string(resp.Header.Peek("Content-Type")), expectedContentType)
}

func setContextValueMiddleware(next app.HandlerFunc, key string, value interface{}) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Set(key, value)
		next(ctx, c)
	}
}

func TestConsumingBodyOnNextConn(t *testing.T) {
	t.Parallel()

	reqNum := 0
	ch := make(chan *http.Request)
	servech := make(chan error, 1)

	opt := config.NewOptions([]config.Option{})
	opt.Addr = "127.0.0.1:10025"
	engine := route.NewEngine(opt)
	handler := func(res http.ResponseWriter, req *http.Request) {
		reqNum++
		ch <- req
	}

	hertzHandler := NewHertzHTTPHandler(http.HandlerFunc(handler))

	engine.POST("/", hertzHandler)
	go engine.Run()
	defer func() {
		engine.Close()
	}()
	time.Sleep(time.Millisecond * 500)

	c, _ := client.NewClient()

	go func() {
		req := protocol.AcquireRequest()
		resp := protocol.AcquireResponse()
		defer func() {
			protocol.ReleaseRequest(req)
			protocol.ReleaseResponse(resp)
		}()
		req.SetRequestURI("http://127.0.0.1:10025")
		req.SetMethod("POST")
		servech <- c.Do(context.Background(), req, resp)
		servech <- c.Do(context.Background(), req, resp)
	}()

	var req *http.Request
	req = <-ch
	if req == nil {
		t.Fatal("Got nil first request.")
	}
	if req.Method != "POST" {
		t.Errorf("For request #1's method, got %q; expected %q",
			req.Method, "POST")
	}

	req = <-ch
	if req == nil {
		t.Fatal("Got nil first request.")
	}
	if req.Method != "POST" {
		t.Errorf("For request #2's method, got %q; expected %q",
			req.Method, "POST")
	}

	if serveerr := <-servech; serveerr != nil {
		t.Errorf("Serve returned %q; expected EOF", serveerr)
	}
}

func TestCopyToHertzRequest(t *testing.T) {
	req := http.Request{
		Method:     "GET",
		RequestURI: "/test",
		URL: &url.URL{
			Scheme: "http",
			Host:   "test.com",
		},
		Proto:  "HTTP/1.1",
		Header: http.Header{},
	}
	req.Header.Set("key1", "value1")
	req.Header.Add("key2", "value2")
	req.Header.Add("key2", "value22")
	hertzReq := protocol.Request{}
	err := adaptor.CopyToHertzRequest(&req, &hertzReq)
	assert.Nil(t, err)
	assert.DeepEqual(t, req.Method, string(hertzReq.Method()))
	assert.DeepEqual(t, req.RequestURI, string(hertzReq.Path()))
	assert.DeepEqual(t, req.Proto, hertzReq.Header.GetProtocol())
	assert.DeepEqual(t, req.Header.Get("key1"), hertzReq.Header.Get("key1"))
	valueSlice := make([]string, 0, 2)
	hertzReq.Header.VisitAllCustomHeader(func(key, value []byte) {
		if strings.ToLower(string(key)) == "key2" {
			valueSlice = append(valueSlice, string(value))
		}
	})

	assert.DeepEqual(t, req.Header.Values("key2"), valueSlice)

	assert.DeepEqual(t, 3, hertzReq.Header.Len())
}

func TestParseArgs(t *testing.T) {
	t.Parallel()

	opt := config.NewOptions([]config.Option{})
	opt.Addr = "127.0.0.1:10026"
	engine := route.NewEngine(opt)
	handler := func(resp http.ResponseWriter, req *http.Request) {
		queryParams := req.URL.Query()

		paramValue := queryParams.Get("test")

		assert.DeepEqual(t, paramValue, "test_value")
	}

	hertzHandler := NewHertzHTTPHandler(http.HandlerFunc(handler))

	engine.GET("/", hertzHandler)
	go engine.Run()
	defer func() {
		engine.Close()
	}()
	time.Sleep(time.Millisecond * 500)

	c, _ := client.NewClient()

	req := protocol.AcquireRequest()
	resp := protocol.AcquireResponse()
	defer func() {
		protocol.ReleaseRequest(req)
		protocol.ReleaseResponse(resp)
	}()
	req.SetRequestURI("http://127.0.0.1:10026/?test=test_value")
	req.SetMethod("GET")

	err := c.Do(context.Background(), req, resp)
	assert.Nil(t, err)
}

func TestCookies(t *testing.T) {
	t.Parallel()

	opt := config.NewOptions([]config.Option{})
	opt.Addr = "127.0.0.1:10027"
	engine := route.NewEngine(opt)
	handler := func(resp http.ResponseWriter, req *http.Request) {
		c, err := req.Cookie("myCookie1")
		assert.Nil(t, err)
		assert.DeepEqual(t, c.Value, "cookieValue1")
		assert.DeepEqual(t, c.HttpOnly, false)
		assert.DeepEqual(t, c.Secure, false)

		c, err = req.Cookie("myCookie2")
		assert.Nil(t, err)
		assert.DeepEqual(t, c.Value, "cookieValue2")

		c.Secure = true
		c.HttpOnly = true
		c.Domain = "google.com"
		c.Expires = time.Now().Add(24 * time.Hour)
		http.SetCookie(resp, c)

		assert.DeepEqual(t, c.HttpOnly, true)
		assert.DeepEqual(t, c.Secure, true)
		assert.DeepEqual(t, c.Domain, "google.com")
		assert.NotEqual(t, c.Expires, nil)
	}

	hertzHandler := NewHertzHTTPHandler(http.HandlerFunc(handler))

	engine.GET("/", hertzHandler)
	go engine.Run()
	defer func() {
		engine.Close()
	}()
	time.Sleep(time.Millisecond * 500)

	c, _ := client.NewClient()

	req := protocol.AcquireRequest()
	resp := protocol.AcquireResponse()
	defer func() {
		protocol.ReleaseRequest(req)
		protocol.ReleaseResponse(resp)
	}()
	req.SetRequestURI("http://127.0.0.1:10027")
	req.SetMethod("GET")
	req.SetCookie("myCookie1", "cookieValue1")
	req.SetCookie("myCookie2", "cookieValue2")

	err := c.Do(context.Background(), req, resp)
	assert.Nil(t, err)
}

func TestHeaders(t *testing.T) {
	t.Parallel()

	opt := config.NewOptions([]config.Option{})
	opt.Addr = "127.0.0.1:10028"
	engine := route.NewEngine(opt)
	handler := func(resp http.ResponseWriter, req *http.Request) {
		k := req.Header.Get("key1")
		assert.DeepEqual(t, k, "value1")
		c := req.Header.Get("cookie")
		assert.DeepEqual(t, c, "cookie=cookie_value")
		assert.DeepEqual(t, req.Header.Get("Content-Type"), "application/form")

		resp.Header().Add("Content-Encoding", "test")
		_, err := resp.Write([]byte("Content-Encoding: test\n"))
		if err != nil {
			panic(err)
		}
	}

	hertzHandler := NewHertzHTTPHandler(http.HandlerFunc(handler))

	engine.GET("/", hertzHandler)
	go engine.Run()
	defer func() {
		engine.Close()
	}()
	time.Sleep(time.Millisecond * 500)

	c, _ := client.NewClient()

	req := protocol.AcquireRequest()
	resp := protocol.AcquireResponse()
	defer func() {
		protocol.ReleaseRequest(req)
		protocol.ReleaseResponse(resp)
	}()
	req.SetRequestURI("http://127.0.0.1:10028")
	req.SetMethod("GET")
	req.Header.Set("key1", "value1")
	req.Header.SetCookie("cookie", "cookie_value")
	req.Header.SetMethod("GET")
	req.Header.SetContentTypeBytes([]byte("application/form"))

	err := c.Do(context.Background(), req, resp)
	assert.Nil(t, err)

	assert.DeepEqual(t, resp.Header.Get("Content-Encoding"), "test")
}

func TestForm(t *testing.T) {
	t.Parallel()

	opt := config.NewOptions([]config.Option{})
	opt.Addr = "127.0.0.1:10029"
	engine := route.NewEngine(opt)
	handler := func(resp http.ResponseWriter, req *http.Request) {
		assert.DeepEqual(t, req.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
		err := req.ParseForm()
		if err != nil {
			return
		}
		assert.DeepEqual(t, req.FormValue("form_data"), "value")
	}

	hertzHandler := NewHertzHTTPHandler(http.HandlerFunc(handler))

	engine.POST("/", hertzHandler)
	go engine.Run()
	defer func() {
		engine.Close()
	}()
	time.Sleep(time.Millisecond * 500)

	c, _ := client.NewClient()

	req := protocol.AcquireRequest()
	resp := protocol.AcquireResponse()
	defer func() {
		protocol.ReleaseRequest(req)
		protocol.ReleaseResponse(resp)
	}()
	req.SetRequestURI("http://127.0.0.1:10029")
	req.SetMethod("POST")

	req.SetFormData(map[string]string{"form_data": "value"})

	err := c.Do(context.Background(), req, resp)
	assert.Nil(t, err)
}

func TestMultiForm(t *testing.T) {
	t.Parallel()

	opt := config.NewOptions([]config.Option{})
	opt.Addr = "127.0.0.1:10030"
	engine := route.NewEngine(opt)
	handler := func(resp http.ResponseWriter, req *http.Request) {
		assert.NotEqual(t, req.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
		err := req.ParseMultipartForm(32 << 20)
		if err != nil {
			return
		}
		assert.DeepEqual(t, req.FormValue("multiform_data"), "value")
	}

	hertzHandler := NewHertzHTTPHandler(http.HandlerFunc(handler))

	engine.POST("/", hertzHandler)
	go engine.Run()
	defer func() {
		engine.Close()
	}()
	time.Sleep(time.Millisecond * 500)

	c, _ := client.NewClient()

	req := protocol.AcquireRequest()
	resp := protocol.AcquireResponse()
	defer func() {
		protocol.ReleaseRequest(req)
		protocol.ReleaseResponse(resp)
	}()
	req.SetRequestURI("http://127.0.0.1:10030")
	req.SetMethod("POST")

	req.SetMultipartFormData(map[string]string{"multiform_data": "value"})

	err := c.Do(context.Background(), req, resp)
	assert.Nil(t, err)
}

//func TestFile(t *testing.T) {
//	t.Parallel()
//
//	opt := config.NewOptions([]config.Option{})
//	opt.Addr = "127.0.0.1:10031"
//	engine := route.NewEngine(opt)
//	handler := func(resp http.ResponseWriter, req *http.Request) {
//		assert.NotEqual(t, req.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
//
//		err := req.ParseMultipartForm(32 << 20)
//		if err != nil {
//			fmt.Println(err)
//			panic(err)
//		}
//
//		file, m, err := req.FormFile("adaptor")
//		if err != nil {
//			fmt.Println(err)
//			panic(err)
//		}
//
//		assert.DeepEqual(t, m.Filename, "handler.go")
//
//		content, err := ioutil.ReadAll(file)
//		assert.Nil(t, err)
//		assert.DeepEqual(t, string(content), "package adaptor\n")
//
//	}
//
//	hertzHandler := NewHertzHTTPHandler(http.HandlerFunc(handler))
//
//	engine.POST("/", hertzHandler)
//	go engine.Run()
//	defer func() {
//		engine.Close()
//	}()
//	time.Sleep(time.Millisecond * 500)
//
//	c, _ := client.NewClient()
//
//	req := protocol.AcquireRequest()
//	resp := protocol.AcquireResponse()
//	defer func() {
//		protocol.ReleaseRequest(req)
//		protocol.ReleaseResponse(resp)
//	}()
//	req.SetRequestURI("http://127.0.0.1:10031")
//	req.SetMethod("POST")
//	req.SetFile("adaptor", "handler.go")
//	fmt.Println(req.MultipartFiles()[0])
//
//	err := c.Do(context.Background(), req, resp)
//	assert.Nil(t, err)
//
//} todo: fix this test
