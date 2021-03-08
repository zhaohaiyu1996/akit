package meta

import (
	"context"
	"github.com/zhaohaiyu1996/akit/registry"
	"google.golang.org/grpc"
)

type RpcMeta struct {
	// Caller name
	Caller string
	// Service name
	ServiceName string
	// called method
	Method string
	// Caller Cluster
	CallerCluster string
	// ServiceC luster
	ServiceCluster string
	// Trace ID
	TraceID string
	// env
	Env string
	// Caller IDC
	CallerIDC string
	// Service IDC
	ServiceIDC string
	// now node
	CurNode *registry.Node
	// History Nodes
	HistoryNodes []*registry.Node
	// All Nodes
	AllNodes []*registry.Node
	// now connect conn
	Conn *grpc.ClientConn
}

type rpcMetaContextKey struct{}

func GetRpcMeta(ctx context.Context) *RpcMeta {
	meta, ok := ctx.Value(rpcMetaContextKey{}).(*RpcMeta)
	if !ok {
		meta = &RpcMeta{}
	}

	return meta
}

func InitRpcMeta(ctx context.Context, service, method, caller string) context.Context {
	meta := &RpcMeta{
		Method:      method,
		ServiceName: service,
		Caller:      caller,
	}
	return context.WithValue(ctx, rpcMetaContextKey{}, meta)
}
