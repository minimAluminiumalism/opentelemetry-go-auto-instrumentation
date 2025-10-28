package ai

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	llmEntryAttributeKey = attribute.Key("gen_ai.is_entry")
	entryStateStore      sync.Map
)

type entryContextKey struct{}

type entryContextValue struct {
	traceID string
	state   *traceEntryState
}

type traceEntryState struct {
	mu              sync.Mutex
	activeSpanCount int
	entryMarked     bool
}

type LLMEntryAttributeExtractor[REQUEST any, RESPONSE any] struct{}

func (e *LLMEntryAttributeExtractor[REQUEST, RESPONSE]) OnStart(attrs []attribute.KeyValue, ctx context.Context, request REQUEST) ([]attribute.KeyValue, context.Context) {
	span := trace.SpanFromContext(ctx)
	spanContext := span.SpanContext()
	if !spanContext.IsValid() || !spanContext.HasTraceID() {
		return attrs, ctx
	}

	traceID := spanContext.TraceID().String()

	stateAny, ok := entryStateStore.Load(traceID)
	var state *traceEntryState
	if !ok {
		state = &traceEntryState{}
		actual, loaded := entryStateStore.LoadOrStore(traceID, state)
		if loaded {
			state = actual.(*traceEntryState)
		}
	} else {
		state = stateAny.(*traceEntryState)
	}

	markEntry := false
	state.mu.Lock()
	state.activeSpanCount++
	if !state.entryMarked {
		state.entryMarked = true
		markEntry = true
	}
	state.mu.Unlock()

	if markEntry {
		attrs = append(attrs, llmEntryAttributeKey.Bool(true))
	}

	ctx = context.WithValue(ctx, entryContextKey{}, entryContextValue{
		traceID: traceID,
		state:   state,
	})

	return attrs, ctx
}

func (e *LLMEntryAttributeExtractor[REQUEST, RESPONSE]) OnEnd(attrs []attribute.KeyValue, ctx context.Context, request REQUEST, response RESPONSE, err error) ([]attribute.KeyValue, context.Context) {
	if value, ok := ctx.Value(entryContextKey{}).(entryContextValue); ok && value.state != nil {
		value.state.mu.Lock()
		value.state.activeSpanCount--
		shouldDelete := value.state.activeSpanCount == 0
		value.state.mu.Unlock()
		if shouldDelete {
			entryStateStore.Delete(value.traceID)
		}
	}
	return attrs, ctx
}
