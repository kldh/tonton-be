package tracing

import (
	"context"
	"errors"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var tp *tracesdk.TracerProvider

type LocalSpanExporter struct {
}

func (e *LocalSpanExporter) ExportSpans(ctx context.Context, spans []tracesdk.ReadOnlySpan) error {
	// nothing
	return nil
}

func (e *LocalSpanExporter) Shutdown(ctx context.Context) error {
	// nothing
	return nil
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

func Shutdown(ctx context.Context) error {
	if tp != nil {
		return tp.Shutdown(ctx)
	}
	return errors.New("nil TraceProvider")
}

// StartSpan starts and returns a context with can be pass into EndSpan to finish a span
func StartSpan(ctx context.Context, name string) context.Context {
	ctx, _ = otel.GetTracerProvider().Tracer("tonton-be").Start(ctx, name)
	return ctx
}

// EndSpan finishes the span that is associated with the given context
func EndSpan(ctx context.Context) {
	if sp := trace.SpanFromContext(ctx); sp != nil {
		sp.End()
	}
}

func TraceID(ctx context.Context) string {
	if sp := trace.SpanFromContext(ctx); sp != nil {
		return sp.SpanContext().TraceID().String()
	}

	return ""
}
