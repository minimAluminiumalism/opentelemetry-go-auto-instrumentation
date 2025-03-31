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

type MessageAttrsGetter[REQUEST any, RESPONSE any] interface {
	GetSystem(request REQUEST) string
	GetDestination(request REQUEST) string
	GetDestinationTemplate(request REQUEST) string
	IsTemporaryDestination(request REQUEST) bool
	isAnonymousDestination(request REQUEST) bool
	GetConversationId(request REQUEST) string
	GetMessageBodySize(request REQUEST) int64
	GetMessageEnvelopSize(request REQUEST) int64
	GetMessageId(request REQUEST, response RESPONSE) string
	GetClientId(request REQUEST) string
	GetBatchMessageCount(request REQUEST, response RESPONSE) int64
	GetMessageHeader(request REQUEST, name string) []string
	GetDestinationPartitionId(request REQUEST) string
}
