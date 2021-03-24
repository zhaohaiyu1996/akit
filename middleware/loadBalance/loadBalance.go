package loadBalance

import (
	"context"
	"github.com/zhaohaiyu1996/akit/aerrors"
	"github.com/zhaohaiyu1996/akit/loadBalance"
	"github.com/zhaohaiyu1996/akit/meta"
	"github.com/zhaohaiyu1996/akit/middleware"
)

// Option is load balance option.
type Option func(*options)

type options struct {
	balancer loadBalance.LoadBalance
}

// set a load balance method,It must be carried out
func WithLimit(balancer loadBalance.LoadBalance) Option {
	return func(o *options) {
		o.balancer = balancer
	}
}

func NewLoadBalanceMiddleware(opts ...Option) middleware.MiddleWare {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	if options.balancer == nil {
		panic("WithLimit method not be carried out")
	}

	return func(next middleware.MiddleWareFunc) middleware.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			//从ctx获取rpc的metadata
			rpcMeta := meta.GetRpcMeta(ctx)
			if len(rpcMeta.AllNodes) == 0 {
				err = aerrors.ErrorByMessage(aerrors.NotHaveInstance, "grpc meta not have a instance caller: %s service: %s", rpcMeta.Caller, rpcMeta.ServiceName)
				return
			}
			//生成loadbalance的上下文,用来过滤已经选择的节点
			ctx = loadBalance.WithBalanceContext(ctx)
			for {
				rpcMeta.CurNode, err = options.balancer.Select(ctx, rpcMeta.AllNodes)
				if err != nil {
					return
				}
				//logs.Debug(ctx, "select node:%#v", rpcMeta.CurNode)
				rpcMeta.HistoryNodes = append(rpcMeta.HistoryNodes, rpcMeta.CurNode)
				resp, err = next(ctx, req)
				if err != nil {
					// if connect error reset
					if aerrors.Is(err, aerrors.NotHaveInstance) {
						continue
					}
					return
				}
				break
			}
			return
		}
	}
}
