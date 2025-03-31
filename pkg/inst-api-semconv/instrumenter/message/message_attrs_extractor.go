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

package message

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
)

type MessageOperation string

const PUBLISH MessageOperation = "publish"
const RECEIVE MessageOperation = "receive"
const PROCESS MessageOperation = "process"

type MessageAttrsExtractor[REQUEST any, RESPONSE any, GETTER MessageAttrsGetter[REQUEST, RESPONSE]] struct {
	getter    GETTER
	operation MessageOperation
}

func (m *MessageAttrsExtractor[REQUEST, RESPONSE, GETTER]) GetSpanKey() attribute.Key {
	switch m.operation {
	case PUBLISH:
		return utils.PRODUCER_KEY
	case RECEIVE:
		return utils.CONSUMER_RECEIVE_KEY
	case PROCESS:
		return utils.CONSUMER_PROCESS_KEY
	}
	panic("Operation" + m.operation + "not supported")
}

func (m *MessageAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) ([]attribute.KeyValue, context.Context) {
	messageAttrSystem := m.getter.GetSystem(request)
	isTemporaryDestination := m.getter.IsTemporaryDestination(request)
	if isTemporaryDestination {
		attributes = append(attributes, attribute.KeyValue{
			Key:   semconv.MessagingDestinationTemporaryKey,
			Value: attribute.BoolValue(true),
		}, attribute.KeyValue{
			Key:   semconv.MessagingDestinationNameKey,
			Value: attribute.StringValue("(temporary)"),
		})
	} else {
		attributes = append(attributes, attribute.KeyValue{
			Key:   semconv.MessagingDestinationNameKey,
			Value: attribute.StringValue(m.getter.GetDestination(request)),
		}, attribute.KeyValue{
			Key:   semconv.MessagingDestinationTemplateKey,
			Value: attribute.StringValue(m.getter.GetDestinationTemplate(request)),
		})
	}
	partitionId := m.getter.GetDestinationPartitionId(request)
	if partitionId != "" {
		attributes = append(attributes, attribute.KeyValue{
			Key:   semconv.MessagingDestinationPartitionIDKey,
			Value: attribute.StringValue(partitionId),
		})
	}
	isAnonymousDestination := m.getter.isAnonymousDestination(request)
	if isAnonymousDestination {
		attributes = append(attributes, attribute.KeyValue{
			Key:   semconv.MessagingDestinationAnonymousKey,
			Value: attribute.BoolValue(true),
		})
	}
	attributes = append(attributes, attribute.KeyValue{
		Key:   semconv.MessagingMessageConversationIDKey,
		Value: attribute.StringValue(m.getter.GetConversationId(request)),
	}, attribute.KeyValue{
		Key:   semconv.MessagingMessageBodySizeKey,
		Value: attribute.Int64Value(m.getter.GetMessageBodySize(request)),
	}, attribute.KeyValue{
		Key:   semconv.MessagingMessageEnvelopeSizeKey,
		Value: attribute.Int64Value(m.getter.GetMessageEnvelopSize(request)),
	}, attribute.KeyValue{
		Key:   semconv.MessagingClientIDKey,
		Value: attribute.StringValue(m.getter.GetClientId(request)),
	}, attribute.KeyValue{
		Key:   semconv.MessagingOperationNameKey,
		Value: attribute.StringValue(string(m.operation)),
	}, attribute.KeyValue{
		Key:   semconv.MessagingSystemKey,
		Value: attribute.StringValue(messageAttrSystem),
	})
	return attributes, parentContext
}

func (m *MessageAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) ([]attribute.KeyValue, context.Context) {
	attributes = append(attributes, attribute.KeyValue{
		Key:   semconv.MessagingMessageIDKey,
		Value: attribute.StringValue(m.getter.GetMessageId(request, response)),
	}, attribute.KeyValue{
		Key:   semconv.MessagingBatchMessageCountKey,
		Value: attribute.Int64Value(m.getter.GetBatchMessageCount(request, response)),
	})
	// TODO: add custom captured headers attributes
	return attributes, context
}
