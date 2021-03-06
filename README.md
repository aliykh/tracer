# Tracer package (opentracing)

## Install jaeger via docker compose with minimum configs
```yaml
version: "3.9"

services:
  jaeger:
    image: jaegertracing/all-in-one:latest
    container_name: jaeger
    ports:
      - "6831:6831/udp"
      - "16686:16686"
    restart: unless-stopped
```


## How to use

```go
package main

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"github.com/aliykh/log"
	openTracer "github.com/aliykh/tracer"
)

var ServiceName = "service_name"

// initTracer - initializes global tracer (opentracing)
func initTracer() {
	
	// initialize logger
	logr := log.NewFactory(log.ZapLogger, "debug")

	// prometheus config
	metricsFactory := prometheus.New().Namespace(metrics.NSOptions{Name: ServiceName, Tags: nil})

	//	initialize jaeger tracer
	tracer, tr := openTracer.InitJaeger(ServiceName, metricsFactory, logr)
	tearDowns = append(tearDowns, tr) // teardown for closing tracer conn

	// Set tracer as global
	opentracing.SetGlobalTracer(tracer)

	return
}
```

## Todo
```text
    # Database wrapper - create a wrapper around sql/Driver package or sqlx package
```