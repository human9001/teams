package logger

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	otelLogSdk "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func getOTLPEndpoint() string {

	return "http://localhost:4318"

}
func InitTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpointURL(getOTLPEndpoint()),
		otlptracehttp.WithURLPath("/v1/traces"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("teams-service"),
			semconv.ServiceVersion("1.0.0"),
			attribute.String("env", "local"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, nil
}

func InitLogger(ctx context.Context) (*otelLogSdk.LoggerProvider, error) {
	exp, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpointURL(getOTLPEndpoint()),
		otlploghttp.WithURLPath("/v1/logs"),
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("teams-service"),
			semconv.ServiceVersion("1.0.0"),
			attribute.String("env", "local"),
		),
	)
	if err != nil {
		return nil, err
	}

	lp := otelLogSdk.NewLoggerProvider(
		otelLogSdk.WithProcessor(otelLogSdk.NewBatchProcessor(exp)),
		otelLogSdk.WithResource(res),
	)

	global.SetLoggerProvider(lp)
	return lp, nil
}

func InitMetrics(ctx context.Context) (*sdkmetric.MeterProvider, error) {
	exp, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpointURL(getOTLPEndpoint()),
		otlpmetrichttp.WithURLPath("/v1/metrics"),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("teams-service"),
			semconv.ServiceVersion("1.0.0"),
			attribute.String("env", "local"),
		),
	)
	if err != nil {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp, sdkmetric.WithInterval(10*time.Second))),
	)

	otel.SetMeterProvider(mp)
	return mp, nil

}
