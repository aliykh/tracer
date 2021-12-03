package tracer

import (
	"fmt"
	"github.com/opentracing/opentracing-go"

	"go.uber.org/zap"

	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/aliykh/log"
)

// InitJaeger -
func InitJaeger(serviceName string, metricsFactory metrics.Factory, logger *log.Factory) (opentracing.Tracer, func()) {

	// Jaeger configuration
	cfg := config.Configuration{
		ServiceName: serviceName, // app name
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	_, err := cfg.FromEnv()
	if err != nil {
		fmt.Printf("jaeger env variables from os env err: %v\n", err.Error())
	}

	// logger for jaeger
	jaegerLogger := jaegerLoggerAdapter{logger: logger.Default()}

	// init jaeger tracer
	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaegerLogger),
		config.Metrics(metricsFactory),
		config.Observer(rpcmetrics.NewObserver(metricsFactory, rpcmetrics.DefaultNameNormalizer)),
	)

	if err != nil {
		logger.Default().Fatal("cannot initialize Jaeger Tracer", zap.Error(err))
	}

	// set tracer as the default tracer of the app
	opentracing.SetGlobalTracer(tracer)

	// teardown for closing the tracer
	tr := func() {
		if err := closer.Close(); err != nil {
			logger.Default().Error("tracer close", zap.Any("err", err.Error()))
			return
		}
		logger.Default().Info("tracer closed [teardown]")
	}

	return tracer, tr
}

type jaegerLoggerAdapter struct {
	logger log.Logger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}
