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

package main

import _ "example/demo/otel_rules"

import _ "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
import _ "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
import _ "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
import _ "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
import _ "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
import _ "go.opentelemetry.io/otel"

import _ "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
import _ "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
import _ "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
import _ "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
import _ "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
import _ "go.opentelemetry.io/otel"

import (
	"example/demo/pkg"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go func() {
		pkg.InitDB()
		pkg.SetupHttp()
	}()

	http.ListenAndServe("0.0.0.0:8080", nil)

	signalCh := make(chan os.Signal, 1)

	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	<-signalCh

	os.Exit(0)
}
