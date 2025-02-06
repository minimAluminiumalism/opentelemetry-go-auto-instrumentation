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
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
)

var grpcServerInstrument = BuildGrpcServerInstrumenter()

// func NewServer(opt ...ServerOption) *Server
func grpcServerOnEnter(call api.CallContext, opts ...grpc.ServerOption) {
	if !grpcEnabler.Enable() {
		return
	}
	h := grpc.StatsHandler(NewServerHandler())
	var opt []grpc.ServerOption
	opt = append(opt, h)
	opt = append(opt, opts...)
	call.SetParam(0, opt)
}

// func NewServer(opt ...ServerOption) *Server
func grpcServerOnExit(call api.CallContext, s *grpc.Server) {
	if !grpcEnabler.Enable() {
		return
	}
	return
}

func NewServerHandler(opts ...Option) stats.Handler {
	h := &serverHandler{
		grpcOtelConfig: newConfig(opts, "server"),
	}

	return h
}

type serverHandler struct {
	*grpcOtelConfig
}

// TagConn can attach some information to the given context.
// stats/opentelemetry/server_metrics.go: func (h *serverStatsHandler) TagConn(ctx context.Context, _ *stats.ConnTagInfo)
func (h *serverHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return ctx
}

// HandleConn processes the Conn stats.
// func (h *serverStatsHandler) HandleConn(context.Context, stats.ConnStats)
func (h *serverHandler) HandleConn(ctx context.Context, info stats.ConnStats) {
}

// TagRPC can attach some information to the given context.
// func (h *serverStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo)
func (h *serverHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	var md metadata.MD
	ctx, md = extract(ctx, h.grpcOtelConfig.Propagators)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	nCtx := grpcServerInstrument.Start(ctx, grpcRequest{
		methodName: info.FullMethodName,
		propagators: &grpcMetadataSupplier{
			metadata: &md,
		},
	})

	gctx := gRPCContext{
		methodName: info.FullMethodName,
	}

	return context.WithValue(nCtx, gRPCContextKey{}, &gctx)
}

// HandleRPC processes the RPC stats.
// func (h *serverStatsHandler) HandleRPC(ctx context.Context, rs stats.RPCStats)
func (h *serverHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	isServer := true
	h.handleRPC(ctx, rs, isServer)
}
