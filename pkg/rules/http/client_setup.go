// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"net/http"
)

// TODO: use a interface to filter
var netHttpFilter = utils.DefaultUrlFilter{}

var netHttpClientInstrumenter = BuildNetHttpClientOtelInstrumenter()

// func (t *Transport) RoundTrip(req *Request) (*Response, error)
func clientOnEnter(call api.CallContext, t *http.Transport, req *http.Request) {
	if netHttpFilter.FilterUrl(req.URL) {
		return
	}
	netHttpRequest := &netHttpRequest{
		method: req.Method,
		url:    req.URL,
		header: req.Header,
		host:   req.Host,
		isTls:  req.TLS != nil,
	}
	netHttpRequest.version = getProtocolVersion(req.ProtoMajor, req.ProtoMinor)
	ctx := netHttpClientInstrumenter.Start(req.Context(), netHttpRequest)
	req = req.WithContext(ctx)
	call.SetParam(1, req)
	data := make(map[string]interface{}, 1)
	data["ctx"] = ctx
	call.SetData(data)
	return
}

// func (t *Transport) RoundTrip(req *Request) (*Response, error)
func clientOnExit(call api.CallContext, res *http.Response, err error) {
	data, ok := call.GetData().(map[string]interface{})
	if !ok || data == nil || data["ctx"] == nil {
		return
	}
	ctx := data["ctx"].(context.Context)
	if res != nil {
		netHttpClientInstrumenter.End(ctx, &netHttpRequest{
			method:  res.Request.Method,
			url:     res.Request.URL,
			header:  res.Request.Header,
			version: getProtocolVersion(res.Request.ProtoMajor, res.Request.ProtoMinor),
			host:    res.Request.Host,
			isTls:   res.Request.TLS != nil,
		}, &netHttpResponse{
			statusCode: res.StatusCode,
			header:     res.Header,
		}, err)
	} else {
		netHttpClientInstrumenter.End(ctx, &netHttpRequest{}, &netHttpResponse{
			statusCode: 500,
		}, err)
	}
}
