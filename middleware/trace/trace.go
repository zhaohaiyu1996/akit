package trace

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	traceLog "github.com/opentracing/opentracing-go/log"
	"github.com/zhaohaiyu1996/akit/meta"
	"github.com/zhaohaiyu1996/akit/middleware"
	"github.com/zhaohaiyu1996/akit/trace"
	"google.golang.org/grpc/metadata"
	"log"
)

// Option is tracing option.
type Option func(*options)

type options struct {
	tracer  opentracing.Tracer
	traceID string
}

// WithTracer sets a custom tracer to be used for this middleware, otherwise the opentracing.GlobalTracer is used.
func WithTracer(tracer opentracing.Tracer) Option {
	return func(o *options) {
		o.tracer = tracer
	}
}

// WithTagId set a traceid
func WithTagId(traceID string) Option {
	return func(o *options) {
		o.traceID = traceID
	}
}

// NewTraceServer is a new server middleware for OpenTracing.
func NewTraceServer(opts ...Option) middleware.MiddleWare {
	options := options{
		tracer: opentracing.GlobalTracer(),
	}
	for _, o := range opts {
		o(&options)
	}
	return func(next middleware.MiddleWareFunc) middleware.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			// get grpc's metadata from context
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				// if not have make it
				md = metadata.Pairs()
			}

			parentSpanContext, err := options.tracer.Extract(opentracing.HTTPHeaders, trace.MetadataTextMap(md))
			if err != nil && err != opentracing.ErrSpanContextNotFound {
				log.Printf("trace extract failed, parsing trace information: %v", err)
			}

			serverMeta := meta.GetServerMeta(ctx)
			// start trace this ,method
			serverSpan := options.tracer.StartSpan(
				serverMeta.Method,
				ext.RPCServerOption(parentSpanContext),
				ext.SpanKindRPCServer,
			)

			serverSpan.SetTag(options.traceID, trace.GetTraceId(ctx))
			ctx = opentracing.ContextWithSpan(ctx, serverSpan)
			resp, err = next(ctx, req)
			// if have aerrors record it
			if err != nil {
				ext.Error.Set(serverSpan, true)
				serverSpan.LogFields(traceLog.String("event", "error"), traceLog.String("message", err.Error()))
			}

			serverSpan.Finish()
			return
		}
	}
}

// NewTraceClient is a new server middleware for OpenTracing.
func NewTraceClient(opts ...Option) middleware.MiddleWare {
	options := options{
		tracer: opentracing.GlobalTracer(),
	}
	for _, o := range opts {
		o(&options)
	}
	return func(next middleware.MiddleWareFunc) middleware.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {

			tracer := opentracing.GlobalTracer()
			var parentSpanCtx opentracing.SpanContext
			if parent := opentracing.SpanFromContext(ctx); parent != nil {
				parentSpanCtx = parent.Context()
			}

			opts := []opentracing.StartSpanOption{
				opentracing.ChildOf(parentSpanCtx),
				ext.SpanKindRPCClient,
				opentracing.Tag{Key: string(ext.Component), Value: "koala_rpc"},
				opentracing.Tag{Key: options.traceID, Value: trace.GetTraceId(ctx)},
			}

			rpcMeta := meta.GetRpcMeta(ctx)
			clientSpan := tracer.StartSpan(rpcMeta.ServiceName, opts...)

			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				md = metadata.Pairs()
			}

			if err := tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, trace.MetadataTextMap(md)); err != nil {
				log.Printf("grpc_opentracing: failed serializing trace information: %v", err)
			}

			ctx = metadata.NewOutgoingContext(ctx, md)
			ctx = metadata.AppendToOutgoingContext(ctx, options.traceID, trace.GetTraceId(ctx))
			ctx = opentracing.ContextWithSpan(ctx, clientSpan)

			resp, err = next(ctx, req)
			// if have aerrors record it
			if err != nil {
				ext.Error.Set(clientSpan, true)
				clientSpan.LogFields(traceLog.String("event", "error"), traceLog.String("message", err.Error()))
			}

			clientSpan.Finish()
			return
		}
	}
}
