package trpc

import (
	"context"

	"trpc.group/trpc-go/trpc-go/codec"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

var trpcClientInstrumenter = BuildTrpcClientInstrumenter()

// func (c *client) Invoke(ctx context.Context, reqBody interface{}, rspBody interface{}, opt ...Option) (err error)
func clientTrpcOnEnter(call api.CallContext, _ interface{}, ctx context.Context, reqBody interface{}, rspBody interface{}, opts interface{}) {
	if !trpcEnabler.Enable() {
		return
	}
	msg := codec.Message(ctx)
	req := trpcReq{
		callerMethod:  msg.CallerMethod(),
		callerService: msg.CallerService(),
		calleeMethod:  msg.CalleeMethod(),
		calleeService: msg.CalleeService(),
		msg:           msg,
	}
	newCtx := trpcClientInstrumenter.Start(context.Background(), req)
	data := make(map[string]interface{}, 3)
	data["ctx"] = newCtx
	data["request"] = req
	data["msg"] = msg
	call.SetData(data)
}

// func (c *client) Invoke(ctx context.Context, reqBody interface{}, rspBody interface{}, opt ...Option) (err error)
func clientTrpcOnExit(call api.CallContext, err error) {
	if !trpcEnabler.Enable() {
		return
	}
	data := call.GetData().(map[string]interface{})
	ctx := data["ctx"].(context.Context)
	request := data["request"].(trpcReq)
	msg := data["msg"].(codec.Msg)
	statusCode := 0
	if msg.ServerRspErr() != nil {
		statusCode = int(msg.ServerRspErr().Code)
	}
	trpcClientInstrumenter.End(ctx, request, trpcRes{
		stausCode: statusCode,
		msg:       msg,
	}, err)
}
