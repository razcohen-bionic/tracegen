package main

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
)

func main() {
  SendBasicTrace(GRPCClient(os.Getenv("TRACE_SERVER_GRPC_ENDPOINT")))
}

func GRPCClient(endpoint string) otlptrace.Client {
	return otlptracegrpc.NewClient(otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(endpoint))
}

func SendBasicTrace(client otlptrace.Client) {
	ctx := context.Background()

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Fatalf("creating opentelemetry trace exporter: %s", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(resource.Default()))
	otel.SetTracerProvider(tracerProvider)
	tracer := otel.GetTracerProvider().Tracer("tracegen", trace.WithInstrumentationVersion("0.0.0"), trace.WithSchemaURL(semconv.SchemaURL))
	ctx, span := tracer.Start(ctx, "basic-trace")
	span.End()
	slog.
		With("originator", "client").
		Info("sending basic trace to server")

	if err = tracerProvider.Shutdown(ctx); err != nil {
		log.Fatalf("error during tracer provider shutdown: %s", err)
	}
}
