// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grpc

import (
	"context"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

// func NewClient(target string, opts ...DialOption) (conn *ClientConn, err error)
func grpcNewClientOnEnter(call api.CallContext, target string, opts ...grpc.DialOption) {
	h := grpc.WithStatsHandler(NewClientNewHandler())
	var opt []grpc.DialOption
	opt = append(opt, h)
	opt = append(opt, opts...)
	call.SetParam(1, opt)
}

// func NewClient(target string, opts ...DialOption) (conn *ClientConn, err error)
func grpcNewClientOnExit(call api.CallContext, cc *grpc.ClientConn, err error) {
	return
}

type clientNewHandler struct {
	*grpcOtelConfig
}

func NewClientNewHandler(opts ...Option) stats.Handler {
	h := &clientNewHandler{
		grpcOtelConfig: newConfig(opts, "client"),
	}

	return h
}

// TagRPC can attach some information to the given context.
// func (h *clientStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context
func (h *clientNewHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	nCtx := grpcClientInstrument.Start(ctx, grpcRequest{
		methodName: info.FullMethodName,
	})
	gctx := gRPCContext{
		methodName: info.FullMethodName,
	}
	return inject(context.WithValue(nCtx, gRPCContextKey{}, &gctx), h.grpcOtelConfig.Propagators, info.FullMethodName)
}

// HandleRPC processes the RPC stats.
func (h *clientNewHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	isServer := false
	h.handleRPC(ctx, rs, isServer)
}

// TagConn can attach some information to the given context.
func (h *clientNewHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return ctx
}

// HandleConn processes the Conn stats.
func (h *clientNewHandler) HandleConn(context.Context, stats.ConnStats) {
	// no-op
}
