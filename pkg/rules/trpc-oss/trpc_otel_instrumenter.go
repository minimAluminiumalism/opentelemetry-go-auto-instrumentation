package trpc

import (
	"fmt"

	"trpc.group/trpc-go/trpc-go/codec"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/rpc"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/version"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/trace"
)

var trpcEnabler = instrumenter.NewDefaultInstrumentEnabler()

type trpcAttrsGetter struct {
}

func (t trpcAttrsGetter) GetSystem(request trpcReq) string {
	return "trpc"
}

func (t trpcAttrsGetter) GetService(request trpcReq) string {
	return request.msg.CallerService()
}

func (t trpcAttrsGetter) GetMethod(request trpcReq) string {
	return request.msg.CallerMethod()
}

func (t trpcAttrsGetter) GetServerAddress(request trpcReq) string {
	return ""
}

type trpcStatusCodeExtractor[REQUEST trpcReq, RESPONSE trpcRes] struct {
}

func (t trpcStatusCodeExtractor[REQUEST, RESPONSE]) Extract(span trace.Span, request trpcReq, response trpcRes, err error) {
	statusCode := response.stausCode
	if statusCode != 0 {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, fmt.Sprintf("trpc error status code %d", statusCode))
		}
	}
}

type trpcRequestCarrier struct {
	reqHeader codec.Msg
}

func (t trpcRequestCarrier) Get(key string) string {
	return string(t.reqHeader.ClientMetaData()[key])
}

func (t trpcRequestCarrier) Set(key string, value string) {
	md := codec.MetaData{}
	md[key] = []byte(value)
	// Need to determine whether the value corresponding to the key exists?
	t.reqHeader.WithClientMetaData(md)
}

func (t trpcRequestCarrier) Keys() []string {
	vals := []string{}
	for _, byteV := range t.reqHeader.ClientMetaData() {
		vals = append(vals, string(byteV))
	}
	return vals
}

func BuildTrpcClientInstrumenter() instrumenter.Instrumenter[trpcReq, trpcRes] {
	builder := instrumenter.Builder[trpcReq, trpcRes]{}
	clientGetter := trpcAttrsGetter{}
	return builder.Init().SetSpanStatusExtractor(&trpcStatusCodeExtractor[trpcReq, trpcRes]{}).SetSpanNameExtractor(&rpc.RpcSpanNameExtractor[trpcReq]{Getter: clientGetter}).
		SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[trpcReq]{}).
		AddAttributesExtractor(&rpc.ClientRpcAttrsExtractor[trpcReq, trpcRes, trpcAttrsGetter]{}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.TRPCGO_CLIENT_SCOPE_NAME,
			Version: version.Tag,
		}).
		AddOperationListeners(rpc.RpcClientMetrics("trpc.client")).
		BuildPropagatingToDownstreamInstrumenter(
			func(n trpcReq) propagation.TextMapCarrier {
				return trpcRequestCarrier{reqHeader: n.msg}
			},
			otel.GetTextMapPropagator(),
		)
}
