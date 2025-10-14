package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
)

func main() {
	tracer := otel.Tracer("test-tracer")
	meter := otel.Meter("test-meter")

	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "test-span")
	defer span.End()

	counter, err := meter.Int64Counter("test.counter")
	if err != nil {
		panic(err)
	}
	counter.Add(ctx, 1)

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req = req.WithContext(ctx)
	client := &http.Client{Timeout: 1 * time.Second}
	_, _ = client.Do(req)

	time.Sleep(100 * time.Millisecond)

	fmt.Println("Multi exporter test completed")
}
