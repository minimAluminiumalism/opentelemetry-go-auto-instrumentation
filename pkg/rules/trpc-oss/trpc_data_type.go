package trpc

import (
	"trpc.group/trpc-go/trpc-go/codec"
)

type trpcReq struct {
	callerMethod  string
	callerService string
	calleeMethod  string
	calleeService string
	msg           codec.Msg
}

type trpcRes struct {
	callerMethod  string
	callerService string
	calleeMethod  string
	calleeService string
	stausCode     int
	msg           codec.Msg
}
