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

package instrumenter

import (
	"context"

	"github.com/alibaba/loongsuite-go-agent/pkg/inst-api/utils"

	"go.opentelemetry.io/otel/attribute"
	ottrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var scopeToCategory = map[string]string{
	// http
	"loongsuite.instrumentation.fasthttp":      "http",
	"loongsuite.instrumentation.nethttp":       "http",
	"loongsuite.instrumentation.hertz":         "http",
	"loongsuite.instrumentation.fiber":         "http",
	"loongsuite.instrumentation.elasticsearch": "http",

	// grpc
	"loongsuite.instrumentation.grpc":   "rpc",
	"loongsuite.instrumentation.trpc":   "rpc",
	"loongsuite.instrumentation.kitex":  "rpc",
	"loongsuite.instrumentation.dubbo":  "rpc",
	"loongsuite.instrumentation.gomicro": "rpc",

	// database
	"loongsuite.instrumentation.databasesql": "db",
	"loongsuite.instrumentation.goredisv9":   "db",
	"loongsuite.instrumentation.goredisv8":   "db",
	"loongsuite.instrumentation.redigo":      "db",
	"loongsuite.instrumentation.mongo":       "db",
	"loongsuite.instrumentation.gorm":        "db",
	"loongsuite.instrumentation.gopg":        "db",
	"loongsuite.instrumentation.gocql":       "db",
	"loongsuite.instrumentation.sqlx":        "db",

	// messaging
	"loongsuite.instrumentation.amqp091":   "messaging",
	"loongsuite.instrumentation.kafka-go":  "messaging",
	"loongsuite.instrumentation.rocketmq":  "messaging",

	// ai/llm
	"loongsuite.instrumentation.eino":      "ai",
	"loongsuite.instrumentation.langchain": "ai",

	// other
	"loongsuite.instrumentation.kratos":        "http",
	"loongsuite.instrumentation.mcp":           "rpc",
	"loongsuite.instrumentation.k8s-client-go": "http",
	"loongsuite.instrumentation.sentinel":      "other",
}

// getScopeKey returns the appropriate span key based on scope name and span kind
func getScopeKey(scopeName string, spanKind trace.SpanKind) attribute.Key {
	category := scopeToCategory[scopeName]
	switch category {
	case "http":
		if spanKind == trace.SpanKindClient {
			return utils.HTTP_CLIENT_KEY
		}
		return utils.HTTP_SERVER_KEY
	case "rpc":
		if spanKind == trace.SpanKindClient {
			return utils.RPC_CLIENT_KEY
		}
		return utils.RPC_SERVER_KEY
	case "db":
		return utils.DB_CLIENT_KEY
	default:
		return ""
	}
}

type SpanSuppressor interface {
	StoreInContext(context context.Context, spanKind trace.SpanKind, span trace.Span) context.Context
	ShouldSuppress(parentContext context.Context, spanKind trace.SpanKind) bool
}

type NoopSpanSuppressor struct {
}

func NewNoopSpanSuppressor() *NoopSpanSuppressor {
	return &NoopSpanSuppressor{}
}

func (n *NoopSpanSuppressor) StoreInContext(context context.Context, spanKind trace.SpanKind, span trace.Span) context.Context {
	return context
}

func (n *NoopSpanSuppressor) ShouldSuppress(parentContext context.Context, spanKind trace.SpanKind) bool {
	return false
}

type SpanKeySuppressor struct {
	spanKeys []attribute.Key
}

func NewSpanKeySuppressor(spanKeys []attribute.Key) *SpanKeySuppressor {
	return &SpanKeySuppressor{spanKeys: spanKeys}
}

func (s *SpanKeySuppressor) StoreInContext(ctx context.Context, spanKind trace.SpanKind, span trace.Span) context.Context {
	// do nothing
	return ctx
}

func (s *SpanKeySuppressor) ShouldSuppress(parentContext context.Context, spanKind trace.SpanKind) bool {
	for _, spanKey := range s.spanKeys {
		span := trace.SpanFromContext(parentContext)
		if s, ok := span.(ottrace.ReadOnlySpan); ok {
			instScopeName := s.InstrumentationScope().Name
			if instScopeName != "" {
				parentSpanKind := s.SpanKind()
				parentSpanKey := getScopeKey(instScopeName, parentSpanKind)
				if spanKey != parentSpanKey {
					return false
				}
			}
		} else {
			return false
		}
	}
	return true
}

func NewSpanKindSuppressor() *SpanKindSuppressor {
	return &SpanKindSuppressor{}
}

func (s *SpanKindSuppressor) StoreInContext(context context.Context, spanKind trace.SpanKind, span trace.Span) context.Context {
	// do nothing
	return context
}

func (s *SpanKindSuppressor) ShouldSuppress(parentContext context.Context, spanKind trace.SpanKind) bool {
	span := trace.SpanFromContext(parentContext)
	if readOnlySpan, ok := span.(ottrace.ReadOnlySpan); ok {
		instScopeName := readOnlySpan.InstrumentationScope().Name
		if instScopeName != "" {
			// Now we compare the actual span kinds directly
			// since scope name no longer distinguishes client/server
			parentSpanKind := readOnlySpan.SpanKind()
			if spanKind != parentSpanKind {
				return false
			}
		}
	} else {
		return false
	}
	return true
}
