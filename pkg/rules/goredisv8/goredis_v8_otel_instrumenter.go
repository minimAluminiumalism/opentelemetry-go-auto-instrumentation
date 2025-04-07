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

package goredisv8

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/db"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/version"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

type goRedisV8AttrsGetter struct {
}

func (d goRedisV8AttrsGetter) GetSystem(request redisv8Data) string {
	return "redis"
}

func (d goRedisV8AttrsGetter) GetServerAddress(request redisv8Data) string {
	return request.Host
}

func (d goRedisV8AttrsGetter) GetDbNamespace(request redisv8Data) string {
	return ""
}

func (d goRedisV8AttrsGetter) GetBatchSize(request redisv8Data) int {
	return 0
}

func (d goRedisV8AttrsGetter) GetStatement(request redisv8Data) string {
	b := make([]byte, 0, 64)

	for i, arg := range request.cmd.Args() {
		if i > 0 {
			b = append(b, ' ')
		}
		b = redisV8AppendArg(b, arg)
	}

	if err := request.cmd.Err(); err != nil && err != redis.Nil {
		b = append(b, ": "...)
		b = append(b, err.Error()...)
	}

	if cmd, ok := request.cmd.(*redis.Cmd); ok {
		b = append(b, ": "...)
		b = redisV8AppendArg(b, cmd)
	}
	return redisV8String(b)
}

func (d goRedisV8AttrsGetter) GetCollection(request redisv8Data) string {
	// TBD: We need to implement retrieving the collection later.
	return ""
}

func (d goRedisV8AttrsGetter) GetOperation(request redisv8Data) string {
	return request.cmd.FullName()
}

func (d goRedisV8AttrsGetter) GetParameters(request redisv8Data) []any {
	return nil
}

func BuildRedisv8Instrumenter() instrumenter.Instrumenter[redisv8Data, any] {
	builder := instrumenter.Builder[redisv8Data, any]{}
	getter := goRedisV8AttrsGetter{}
	return builder.Init().SetSpanNameExtractor(&db.DBSpanNameExtractor[redisv8Data]{Getter: getter}).SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[redisv8Data]{}).
		AddAttributesExtractor(&db.DbClientAttrsExtractor[redisv8Data, any, db.DbClientAttrsGetter[redisv8Data]]{Base: db.DbClientCommonAttrsExtractor[redisv8Data, any, db.DbClientAttrsGetter[redisv8Data]]{Getter: getter}}).
		AddOperationListeners(db.DbClientMetrics("nosql.goredisv8")).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.GO_REDIS_V8_SCOPE_NAME,
			Version: version.Tag,
		}).
		BuildInstrumenter()
}
