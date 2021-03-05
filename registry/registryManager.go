package registry

import (
	"context"
	"fmt"
	"sync"
)

// init PluginMgr
var pluginMgr = &PluginMgr{
	plugins: make(map[string]Registry),
}

// registry Plugin manager struct
type PluginMgr struct {
	plugins map[string]Registry
	lock    sync.Mutex
}

// registerPlugin is registry a service registry plugin
func (p *PluginMgr) registerPlugin(plugin Registry) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	_, ok := p.plugins[plugin.Name()]
	if ok {
		return fmt.Errorf("duplicate registry plugin")
	}

	p.plugins[plugin.Name()] = plugin
	return nil
}

// initRegistry is init service registry
func (p *PluginMgr) initRegistry(ctx context.Context, name string,
	opts ...Option) (Registry, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	plugin, ok := p.plugins[name]
	if !ok {
		return nil, fmt.Errorf("plugin %s not exists", name)
	}
	err := plugin.Init(ctx, opts...)
	return plugin, err
}

// RegisterPlugin is register a service registry plugin
func RegisterPlugin(registry Registry)  error {
	return pluginMgr.registerPlugin(registry)
}

// init Registry
func InitRegistry(ctx context.Context, name string, opts ...Option) (Registry, error) {
	return pluginMgr.initRegistry(ctx, name, opts...)
}
