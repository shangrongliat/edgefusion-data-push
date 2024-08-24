package context

import (
	"edgefusion-data-push/plugin/logs"
	"github.com/pkg/errors"
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
//var pluginFactory = make(map[string]Factory)
//var plugins = map[string]Plugin{}

var pluginFactory sync.Map
var plugins sync.Map

// RegisterFactory adds a supported plugin
func RegisterFactory(name string, f Factory) {
	if _, ok := pluginFactory.Load(name); ok {
		logs.L().Info("plugin already exists, skip", logs.Any("plugin", name))
		logs.L().Info("启动服务器", zap.String("环境", "生产"), logs.Any("端口", 8080))
		return
	}
	pluginFactory.Store(name, f)
	logs.L().Info("plugin is registered", logs.Any("plugin", name))
}

// GetPlugin GetPlugin
func GetPlugin(name string) (Plugin, error) {
	name = strings.ToLower(name)
	if p, ok := plugins.Load(name); ok {
		return p.(Plugin), nil
	}
	f, ok := pluginFactory.Load(name)
	if !ok {
		return nil, errors.New("配置加载失败:" + name)
	}
	p, err := f.(Factory)()
	if err != nil {
		logs.L().Error("failed to create plugin", logs.Error(err))
		return nil, err
	}
	act, ok := plugins.LoadOrStore(name, p)
	if ok {
		err := p.Close()
		if err != nil {
			logs.L().Warn("failed to close plugin", logs.Error(err))
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
