package server

import (
	"context"
	"fmt"
	"io/ioutil"
	v1 "rag/api/gateway/v1"
	"rag/app/gateway/internal/conf"
	"rag/app/gateway/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	md "github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, gateway *service.GatewayService, logger log.Logger) *http.Server {
	exporter, err := stdouttrace.New(stdouttrace.WithWriter(ioutil.Discard))
	if err != nil {
		fmt.Printf("creating stdout exporter: %v", err)
		panic(err)
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String("gateway")),
		))
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	opts = append(opts, http.Middleware(
		metadata.Server(),
		metadata.Client(),
		tracing.Server(tracing.WithTracerProvider(tp)),
		func(h middleware.Handler) middleware.Handler {
			return func(ctx context.Context, req interface{}) (interface{}, error) {
				a, b := md.FromServerContext(ctx)
				if b {
					a.Range(func(k string, v []string) bool {
						fmt.Println(k, v)
						return true
					})
				}
				tr, _ := transport.FromServerContext(ctx)

				if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
					tr.ReplyHeader().Set("X-Trace-Id", span.TraceID().String())
				}
				return h(ctx, req)
			}
		},
		// validate.Validator(),
	))

	srv := http.NewServer(opts...)
	v1.RegisterGatewayHTTPServer(srv, gateway)

	return srv
}
