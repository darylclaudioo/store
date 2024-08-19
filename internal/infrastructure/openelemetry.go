package infrastructure

import (
	"context"
	"errors"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	otelTrace "go.opentelemetry.io/otel/trace"
)

func SetupOTelSDK(ctx context.Context, serviceName, serviceVersion, endpoint, env string) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	handleErr := func(inErr error) {
	}

	res, err := newResource(serviceName, serviceVersion, env)
	if err != nil {
		handleErr(err)
		return
	}

	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	tracerProvider, err := newTraceProvider(ctx, endpoint, res)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	meterProvider, err := newMeterProvider(ctx, endpoint, res)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return
}

func newResource(serviceName, serviceVersion, env string) (*resource.Resource, error) {
	extraResources, _ := resource.New(
		context.Background(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithContainer(),
		resource.WithHost(),
		resource.WithAttributes(semconv.ServiceName(serviceName), semconv.ServiceVersion(serviceVersion), semconv.DeploymentEnvironment(env)),
	)
	return resource.Merge(resource.Default(), extraResources)
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(ctx context.Context, endpoint string, res *resource.Resource) (*trace.TracerProvider, error) {
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(endpoint))
	if err != nil {
		log.Fatalf("new otlp trace grpc exporter failed: %v", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(traceProvider)

	return traceProvider, nil
}

func newMeterProvider(_ context.Context, _ string, res *resource.Resource) (*metric.MeterProvider, error) {
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithView(metric.NewView(
			metric.Instrument{Scope: instrumentation.Scope{Name: "go.opentelemetry.io/contrib/google.golang.org/grpc/otelgrpc"}},
			metric.Stream{Aggregation: metric.AggregationDrop{}},
		)),
		metric.WithReader(metric.NewPeriodicReader(metricExporter)))
	otel.SetMeterProvider(meterProvider)
	return meterProvider, nil
}

func Tracer() otelTrace.Tracer {
	return otel.Tracer("store-service")
}
