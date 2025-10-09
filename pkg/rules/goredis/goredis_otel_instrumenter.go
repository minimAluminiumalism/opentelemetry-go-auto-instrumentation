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

package goredis

import (
	"fmt"
	"github.com/alibaba/loongsuite-go-agent/pkg/inst-api-semconv/instrumenter/db"
	"github.com/alibaba/loongsuite-go-agent/pkg/inst-api/instrumenter"
	"github.com/alibaba/loongsuite-go-agent/pkg/inst-api/utils"
	"github.com/alibaba/loongsuite-go-agent/pkg/inst-api/version"
	redis "github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"strconv"
	"time"
	"unicode/utf8"
	"unsafe"
)

type goRedisAttrsGetter struct {
}

func (d goRedisAttrsGetter) GetSystem(request goRedisRequest) string {
	return "redis"
}

func (d goRedisAttrsGetter) GetServerAddress(request goRedisRequest) string {
	return request.endpoint
}

func (d goRedisAttrsGetter) GetStatement(request goRedisRequest) string {
	b := make([]byte, 0, 64)

	for i, arg := range request.cmd.Args() {
		if i > 0 {
			b = append(b, ' ')
		}
		b = redisV9AppendArg(b, arg)
	}

	if err := request.cmd.Err(); err != nil {
		b = append(b, ": "...)
		b = append(b, err.Error()...)
	}

	if cmd, ok := request.cmd.(*redis.Cmd); ok {
		b = append(b, ": "...)
		b = redisV9AppendArg(b, cmd)
	}

	return redisV9String(b)
}

func (d goRedisAttrsGetter) GetOperation(request goRedisRequest) string {
	return request.cmd.FullName()
}

func (d goRedisAttrsGetter) GetCollection(request goRedisRequest) string {
	// TBD: We need to implement retrieving the collection later.
	return ""
}

func (d goRedisAttrsGetter) GetParameters(request goRedisRequest) []any {
	return nil
}

func (d goRedisAttrsGetter) GetDbNamespace(request goRedisRequest) string {
	return ""
}

func (d goRedisAttrsGetter) GetBatchSize(request goRedisRequest) int {
	return 0
}

func BuildGoRedisOtelInstrumenter() instrumenter.Instrumenter[goRedisRequest, any] {
	builder := instrumenter.Builder[goRedisRequest, any]{}
	getter := goRedisAttrsGetter{}
	return builder.Init().SetSpanNameExtractor(&db.DBSpanNameExtractor[goRedisRequest]{Getter: getter}).SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[goRedisRequest]{}).
		AddAttributesExtractor(&db.DbClientAttrsExtractor[goRedisRequest, any, db.DbClientAttrsGetter[goRedisRequest]]{Base: db.DbClientCommonAttrsExtractor[goRedisRequest, any, db.DbClientAttrsGetter[goRedisRequest]]{Getter: getter}}).
		AddOperationListeners(db.DbClientMetrics("nosql.goredisv9")).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.GO_REDIS_V9_SCOPE_NAME,
			Version: version.Tag,
		}).
		BuildInstrumenter()
}

func redisV9String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func redisV9AppendUTF8String(dst []byte, src []byte) []byte {
	dst = append(dst, src...)
	return dst
}

func redisV9Bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func redisV9AppendArg(b []byte, v interface{}) []byte {
	switch v := v.(type) {
	case nil:
		return append(b, "<nil>"...)
	case string:
		bts := redisV9Bytes(v)
		if utf8.Valid(bts) {
			return redisV9AppendUTF8String(b, bts)
		} else {
			return redisV9AppendUTF8String(b, redisV9Bytes("<string>"))
		}
	case []byte:
		if utf8.Valid(v) {
			return redisV9AppendUTF8String(b, v)
		} else {
			return redisV9AppendUTF8String(b, redisV9Bytes("<byte>"))
		}
	case int:
		return strconv.AppendInt(b, int64(v), 10)
	case int8:
		return strconv.AppendInt(b, int64(v), 10)
	case int16:
		return strconv.AppendInt(b, int64(v), 10)
	case int32:
		return strconv.AppendInt(b, int64(v), 10)
	case int64:
		return strconv.AppendInt(b, v, 10)
	case uint:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint8:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint16:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint32:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint64:
		return strconv.AppendUint(b, v, 10)
	case float32:
		return strconv.AppendFloat(b, float64(v), 'f', -1, 64)
	case float64:
		return strconv.AppendFloat(b, v, 'f', -1, 64)
	case bool:
		if v {
			return append(b, "true"...)
		}
		return append(b, "false"...)
	case time.Time:
		return v.AppendFormat(b, time.RFC3339Nano)
	default:
		return append(b, fmt.Sprint(v)...)
	}
}
