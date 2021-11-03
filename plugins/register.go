package plugins

import (
	"fmt"
	log "github.com/fhmq/hmq/logger"
	"go.uber.org/zap"
	"io"
	"strings"
	"sync"
)

// Plugin interfaces
type Plugin interface {
	io.Closer
}

// Factory create engine by given config
type Factory func() (Plugin, error)

// PluginFactory contains all supported plugin factory
var pluginFactory sync.Map
var plugins sync.Map

// RegisterFactory adds a supported plugin
func RegisterFactory(name string, f Factory) {
	if _, ok := pluginFactory.Load(name); ok {
		log.Info("plugin already exists, skip", zap.String("plugin", name))
		return
	}
	pluginFactory.Store(name, f)
	log.Info("plugin is registered", zap.String("plugin", name))
}

// GetPlugin GetPlugin
func GetPlugin(name string) (Plugin, error) {
	name = strings.ToLower(name)
	if p, ok := plugins.Load(name); ok {
		return p.(Plugin), nil
	}
	f, ok := pluginFactory.Load(name)
	if !ok {
		return nil, fmt.Errorf("plugin [%s] not found", name)
	}
	p, err := f.(Factory)()
	if err != nil {
		log.Error("failed to create plugin", zap.Error(err))
		return nil, err
	}
	act, ok := plugins.LoadOrStore(name, p)
	if ok {
		err := p.Close()
		if err != nil {
			log.Warn("failed to close plugin", zap.Error(err))
		}
		return act.(Plugin), nil
	}
	return p, nil
}

// ClosePlugins ClosePlugins
func ClosePlugins() {
	plugins.Range(func(key, value interface{}) bool {
		value.(Plugin).Close()
		return true
	})
}
