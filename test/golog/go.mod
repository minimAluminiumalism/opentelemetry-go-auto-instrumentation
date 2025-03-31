module golog

go 1.23

toolchain go1.23.3

replace github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier => ../../../opentelemetry-go-auto-instrumentation/test/verifier

replace github.com/alibaba/opentelemetry-go-auto-instrumentation => ../../../opentelemetry-go-auto-instrumentation

require github.com/alibaba/opentelemetry-go-auto-instrumentation v0.0.0-00010101000000-000000000000

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/sdk v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
)
