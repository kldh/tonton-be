package tracing

import (
	"context"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type LocalSpanExporter struct {
}

func (e *LocalSpanExporter) ExportSpans(ctx context.Context, spans []tracesdk.ReadOnlySpan) error {
	// nothing
}

func (e *LocalSpanExporter) Shutdown(ctx context.Context) error {
	// nothing
}

func Init() {
	provider := &LocalSpanExporter{}
	initTracer(provider, "tonton-be")
}

func initTracer(exp tracesdk.SpanExporter, serviceName string) {
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(b3.New(b3.WithInjectEncoding(b3.B3SingleHeader | b3.B3MultipleHeader)))
}
