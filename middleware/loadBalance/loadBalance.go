package loadBalance

import (
	"context"
	"github.com/zhaohaiyu1996/akit/errors"
	"github.com/zhaohaiyu1996/akit/loadBalance"
	"github.com/zhaohaiyu1996/akit/meta"
	"github.com/zhaohaiyu1996/akit/middleware"
)

func NewLoadBalanceMiddleware(balancer loadBalance.LoadBalance) middleware.MiddleWare {
	return func(next middleware.MiddleWareFunc) middleware.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			//从ctx获取rpc的metadata
			rpcMeta := meta.GetRpcMeta(ctx)
			if len(rpcMeta.AllNodes) == 0 {
				err = errors.ErrorByMessage(errors.NotHaveInstance, "grpc meta not have a instance caller: %s service: %s", rpcMeta.Caller, rpcMeta.ServiceName)
				return
			}
			//生成loadbalance的上下文,用来过滤已经选择的节点
			ctx = loadBalance.WithBalanceContext(ctx)
			for {
				rpcMeta.CurNode, err = balancer.Select(ctx, rpcMeta.AllNodes)
				if err != nil {
					return
				}
				//logs.Debug(ctx, "select node:%#v", rpcMeta.CurNode)
				rpcMeta.HistoryNodes = append(rpcMeta.HistoryNodes, rpcMeta.CurNode)
				resp, err = next(ctx, req)
				if err != nil {
					// if connect error reset
					if errors.Is(err,errors.NotHaveInstance) {
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
