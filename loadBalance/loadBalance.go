package loadBalance

import (
	"context"
	"github.com/zhaohaiyu1996/akit/registry"
)

type LoadBalance interface {
	Name() string
	Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error)
}