module github.com/alibaba/loongsuite-go-agent/pkg/rules/dubbo

go 1.23.0

replace github.com/alibaba/loongsuite-go-agent/pkg => ../../../pkg

require (
	dubbo.apache.org/dubbo-go/v3 v3.3.0
	github.com/alibaba/loongsuite-go-agent/pkg v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.35.0
	go.opentelemetry.io/otel/sdk v1.35.0
	go.opentelemetry.io/otel/trace v1.35.0
)

require (
	github.com/Workiva/go-datastructures v1.0.52 // indirect
	github.com/creasty/defaults v1.5.2 // indirect
	github.com/dubbogo/gost v1.14.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/copier v0.3.5 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
)
