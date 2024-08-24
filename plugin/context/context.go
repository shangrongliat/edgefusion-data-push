package context

import (
	"edgefusion-data-push/plugin/config"
	"edgefusion-data-push/plugin/logs"
	"edgefusion-data-push/plugin/utils"
	"errors"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	ErrSystemCertInvalid  = errors.New("system certificate is invalid")
	ErrSystemCertNotFound = errors.New("system certificate is not found")
)

// Context of service
type Context interface {
	// NodeName returns node name from data.
	NodeName() string
	// AppName returns app name from data.
	AppName() string
	// AppVersion returns application version from data.
	AppVersion() string
	// ServiceName returns service name from data.
	ServiceName() string
	// ConfFile returns config file from data.
	ConfFile() string

	// SystemConfig returns the config of baetyl system from data.
	SystemConfig() *config.SystemConfig

	// Log returns logger interface.
	Log() *zap.Logger

	// Wait waits until exit, receiving SIGTERM and SIGINT signals.
	Wait()
	// WaitChan returns wait channel.
	WaitChan() <-chan os.Signal

	// Load returns the value stored in the map for a key, or nil if no value is present.
	// The ok result indicates whether value was found in the map.
	Load(key interface{}) (value interface{}, ok bool)
	// Store sets the value for a key.
	Store(key, value interface{})
	// LoadOrStore returns the existing value for the key if present.
	// Otherwise, it stores and returns the given value.
	// The loaded result is true if the value was loaded, false if stored.
	LoadOrStore(key, value interface{}) (actual interface{}, loaded bool)
	// Delete deletes the value for a key.
	Delete(key interface{})

	// CheckSystemCert checks system certificate, if certificate is not found or invalid, returns an error.
	CheckSystemCert() error
	// LoadCustomConfig loads custom config.
	// If 'files' is empty, will load config from default path,
	// else the first file path will be used to load config from.
	LoadCustomConfig(cfg interface{}, files ...string) error
}

type ctx struct {
	sync.Map // global cache
	log      *logs.Logger
}

// NewContext creates a new context
func NewContext(confFile string) Context {
	if confFile == "" {
		confFile = os.Getenv(KeyConfFile)
	}

	c := &ctx{}
	c.Store(KeyConfFile, confFile)
	c.Store(KeyNodeName, os.Getenv(KeyNodeName))
	c.Store(KeyAppName, os.Getenv(KeyAppName))
	c.Store(KeyAppVersion, os.Getenv(KeyAppVersion))
	c.Store(KeySvcName, os.Getenv(KeySvcName))

	var lfs []logs.Field
	if c.NodeName() != "" {
		lfs = append(lfs, zap.Any("node", c.NodeName()))
	}
	if c.AppName() != "" {
		lfs = append(lfs, zap.Any("app", c.AppName()))
	}
	if c.ServiceName() != "" {
		lfs = append(lfs, zap.Any("service", c.ServiceName()))
	}
	c.log = logs.L().With(lfs...)
	c.log.Info("to load config file", zap.Any("file", c.ConfFile()))

	sc := &config.SystemConfig{}
	err := c.LoadCustomConfig(sc)
	if err != nil {
		c.log.Error("failed to load system config, to use default config", zap.Error(err))
		config.UnmarshalYAML(nil, sc)
	}
	// populate configuration
	// if not set in config file, to use value from env.
	// if not set in env, to use default value.

	_log, err := logs.Init(sc.Logger, lfs...)
	if err != nil {
		c.log.Error("failed to init logger", logs.Error(err))
	}
	c.log = _log
	c.log.Debug("context is created", zap.Any("file", confFile))
	return c
}

func (c *ctx) NodeName() string {
	v, ok := c.Load(KeyNodeName)
	if !ok {
		return ""
	}
	return v.(string)
}

func (c *ctx) AppName() string {
	v, ok := c.Load(KeyAppName)
	if !ok {
		return ""
	}
	return v.(string)
}

func (c *ctx) AppVersion() string {
	v, ok := c.Load(KeyAppVersion)
	if !ok {
		return ""
	}
	return v.(string)
}

func (c *ctx) ServiceName() string {
	v, ok := c.Load(KeySvcName)
	if !ok {
		return ""
	}
	return v.(string)
}

func (c *ctx) ConfFile() string {
	v, ok := c.Load(KeyConfFile)
	if !ok {
		return ""
	}
	return v.(string)
}

func (c *ctx) SystemConfig() *config.SystemConfig {
	v, ok := c.Load(KeySysConf)
	if !ok {
		return nil
	}
	return v.(*config.SystemConfig)
}

func (c *ctx) Log() *zap.Logger {
	return c.log
}

func (c *ctx) Wait() {
	<-c.WaitChan()
}

func (c *ctx) WaitChan() <-chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	signal.Ignore(syscall.SIGPIPE)
	return sig
}

func (c *ctx) CheckSystemCert() error {
	cfg := c.SystemConfig().Certificate
	if !utils.FileExists(cfg.CA) || !utils.FileExists(cfg.Key) || !utils.FileExists(cfg.Cert) {
		return ErrSystemCertNotFound
	}
	crt, err := ioutil.ReadFile(cfg.Cert)
	if err != nil {
		return err
	}
	info, err := utils.ParseCertificates(crt)
	if err != nil {
		return err
	}
	if len(info) != 1 || len(info[0].Subject.OrganizationalUnit) != 1 ||
		info[0].Subject.OrganizationalUnit[0] != KeyEdgeFusion {
		return ErrSystemCertInvalid
	}
	return nil
}

func (c *ctx) LoadCustomConfig(cfg interface{}, files ...string) error {
	f := c.ConfFile()
	if len(files) > 0 && len(files[0]) > 0 {
		f = files[0]
	}
	if utils.FileExists(f) {
		return config.LoadYAML(f, cfg)
	}
	return config.UnmarshalYAML(nil, cfg)
}
