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

package trpc

import (
	"context"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/errs"
)

var trpcClientInstrumenter = BuildTrpcClientInstrumenter()

// func (c *clientTransport) RoundTrip(ctx context.Context, req []byte, roundTripOpts ...RoundTripOption) (rsp []byte, err error)
// https://github.com/trpc-group/trpc-go/blob/e025145c92d41417fb71574fb486441e629804ac/transport/client_transport.go#L65
func clientTrpcOnEnter(call api.CallContext, _ interface{}, ctx context.Context, request []byte, roundTripOpts interface{}) {
	if !trpcEnabler.Enable() {
		return
	}
	msg := codec.Message(ctx)
	req := trpcReq{
		msg: msg,
	}
	newCtx := trpcClientInstrumenter.Start(context.Background(), req)
	data := make(map[string]interface{}, 3)
	data["ctx"] = newCtx
	data["request"] = req
	call.SetData(data)
}

// func (c *clientTransport) RoundTrip(ctx context.Context, req []byte, roundTripOpts ...RoundTripOption) (rsp []byte, err error)
func clientTrpcOnExit(call api.CallContext, rsp []byte, err error) {
	if !trpcEnabler.Enable() {
		return
	}
	data := call.GetData().(map[string]interface{})
	ctx := data["ctx"].(context.Context)
	request := data["request"].(trpcReq)
	statusCode := 0
	if err != nil {
		// ref: https://github.com/trpc-group/trpc/blob/main/trpc/trpc.proto
		statusCode = int(err.(*errs.Error).Code)
	}
	trpcClientInstrumenter.End(ctx, request, trpcRes{
		stausCode: statusCode,
	}, err)
}
