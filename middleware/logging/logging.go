package middleware

import (
	"context"
	"fmt"
	"github.com/zhaohaiyu1996/akit/log"
	"github.com/zhaohaiyu1996/akit/meta"
	"github.com/zhaohaiyu1996/akit/middleware"
	"time"

	"google.golang.org/grpc/status"
)

type Option func(*options)

type options struct {
	logger log.Logger
}

// WithLogger with middleware logger.
func WithLogger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func NewServerLogMiddleware(opts ...Option) middleware.MiddleWare {
	options := options{
		logger: log.DefaultLogger,
	}
	for _, o := range opts {
		o(&options)
	}
	log := log.NewHelper("middleware/logging", options.logger)
	return func(next middleware.MiddleWareFunc) middleware.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			startTime := time.Now()
			resp, err = next(ctx, req)
			serverMeta := meta.GetServerMeta(ctx)
			errStatus, _ := status.FromError(err)
			cost := time.Since(startTime).Nanoseconds() / 1000
			log.Infof("cost_us:%d", cost)
			log.Infof("method:%s", serverMeta.Method)

			log.Infof("cluster:%s", serverMeta.Cluster)
			log.Infof("env:%s", serverMeta.Env)
			log.Infof("server_ip:%s", serverMeta.ServerIP)
			log.Infof("client_ip:%s", serverMeta.ClientIP)
			log.Infof("idc:%s", serverMeta.IDC)
			log.Infof("result=%v", errStatus.Code())

			return
		}
	}
}

func NewClientLogMiddleware(opts ...Option) middleware.MiddleWare {
	options := options{
		logger: log.DefaultLogger,
	}
	for _, o := range opts {
		o(&options)
	}
	log := log.NewHelper("middleware/logging", options.logger)
	return func(next middleware.MiddleWareFunc) middleware.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			startTime := time.Now()
			resp, err = next(ctx, req)

			rpcMeta := meta.GetRpcMeta(ctx)
			errStatus, _ := status.FromError(err)

			cost := time.Since(startTime).Nanoseconds() / 1000
			log.Infof("cost_us", cost)
			log.Infof("method", rpcMeta.Method)
			log.Infof("server", rpcMeta.ServiceName)

			log.Infof("caller_cluster", rpcMeta.CallerCluster)
			log.Infof("upstream_cluster", rpcMeta.ServiceCluster)
			log.Infof("rpc", 1)
			log.Infof("env", rpcMeta.Env)

			var upstreamInfo string
			for _, node := range rpcMeta.HistoryNodes {
				upstreamInfo += fmt.Sprintf("%s:%d,", node.IP, node.Port)
			}

			log.Infof("upstream", upstreamInfo)
			log.Infof("caller_idc", rpcMeta.CallerIDC)
			log.Infof("upstream_idc", rpcMeta.ServiceIDC)
			log.Infof("result=%v", errStatus.Code())

			return
		}
	}
}
