package config

import (
	"time"

	log "edgefusion-data-push/plugin/logs"
)

// Config config
type Config struct {
	Server  Server     `yaml:"server" json:"server" `
	LogInfo log.Config `yaml:"logger" json:"logger"`
	Task    Task       `yaml:"task" json:"task"`
	Lock    Lock       `yaml:"lock" json:"lock"`
	Cache   struct {
		ExpirationDuration time.Duration `yaml:"expirationDuration" json:"expirationDuration" default:"10m"`
	} `yaml:"cache" json:"cache"`
	Database DatabaseConfig `yaml:"database" json:"database"`
	Minio    MinioConfig    `yaml:"minio" json:"minio"`
	Hook     LiveHook       `yaml:"hook" json:"hook"`
	Mqtt     MqttConfig     `yaml:"mqtt" json:"mqtt"`
}

type CronJob struct {
	CronName string `yaml:"cronName" json:"cronName"`
	CronGap  string `yaml:"cronGap" json:"cronGap" default:"20s"`
}

// Server server config
type Server struct {
	Port         string        `yaml:"port" json:"port"`
	ReadTimeout  time.Duration `yaml:"readTimeout" json:"readTimeout" default:"30s"`
	WriteTimeout time.Duration `yaml:"writeTimeout" json:"writeTimeout" default:"30s"`
	ShutdownTime time.Duration `yaml:"shutdownTime" json:"shutdownTime" default:"3s"`
	Certificate  Certificate   `yaml:",inline" json:",inline"`
}

type Task struct {
	BatchNum        int32 `yaml:"batchNum" json:"batchNum" default:"100"`
	LockExpiredTime int32 `yaml:"lockExpiredTime" json:"lockExpiredTime" default:"60" unit:"second"`
	ScheduleTime    int32 `yaml:"scheduletime" json:"scheduletime" default:"30" unit:"second"`
	ConcurrentNum   int32 `yaml:"concurrentNum" json:"concurrentNum" default:"10"`
	QueueLength     int32 `yaml:"queueLength" json:"queueLength" default:"100"`
}

type Lock struct {
	ExpireTime int64 `yaml:"expireTime" json:"expireTime" default:"5" unit:"second"`
}

type DatabaseConfig struct {
	Type string `yaml:"type" json:"type"` // mysql, postgres, sqlite3
	Url  string `yaml:"url" json:"url"`
}

// ClientConfig client config
type ClientConfig struct {
	Address               string        `yaml:"address" json:"address"`
	Timeout               time.Duration `yaml:"timeout" json:"timeout" default:"30s"`
	KeepAlive             time.Duration `yaml:"keepalive" json:"keepalive" default:"30s"`
	MaxIdleConns          int           `yaml:"maxIdleConns" json:"maxIdleConns" default:"100"`
	IdleConnTimeout       time.Duration `yaml:"idleConnTimeout" json:"idleConnTimeout" default:"90s"`
	TLSHandshakeTimeout   time.Duration `yaml:"tlsHandshakeTimeout" json:"tlsHandshakeTimeout" default:"10s"`
	ExpectContinueTimeout time.Duration `yaml:"expectContinueTimeout" json:"expectContinueTimeout" default:"1s"`
	ByteUnit              string        `yaml:"byteUnit" json:"byteUnit" default:"KB"`
	SpeedLimit            int           `yaml:"speedLimit" json:"speedLimit" default:"0"`
	SyncMaxConcurrency    int           `yaml:"syncMaxConcurrency" json:"syncMaxConcurrency" default:"0"`
	Certificate           `yaml:",inline" json:",inline"`
}

type MqttConfig struct {
	Address              string        `yaml:"address" json:"address"`
	Username             string        `yaml:"username" json:"username"`
	Password             string        `yaml:"password" json:"password"`
	ClientID             string        `yaml:"clientid" json:"clientid"`
	CleanSession         bool          `yaml:"cleansession" json:"cleansession"`
	Timeout              time.Duration `yaml:"timeout" json:"timeout" default:"30s"`
	KeepAlive            time.Duration `yaml:"keepalive" json:"keepalive" default:"30s"`
	MaxReconnectInterval time.Duration `yaml:"maxReconnectInterval" json:"maxReconnectInterval" default:"3m"`
	MaxCacheMessages     int           `yaml:"maxCacheMessages" json:"maxCacheMessages" default:"10"`
	DisableAutoAck       bool          `yaml:"disableAutoAck" json:"disableAutoAck"`
	Subscriptions        []QOSTopic    `yaml:"subscriptions" json:"subscriptions" default:"[]"`
	Certificate          `yaml:",inline" json:",inline"`
}

// QOSTopic topic and qos
type QOSTopic struct {
	QOS   uint32 `yaml:"qos" json:"qos" binding:"min=0,max=1"`
	Topic string `yaml:"topic" json:"topic" binding:"nonzero"`
}

type LoggerConfig struct {
	Level    string `yaml:"level" json:"level" default:"info"`
	Filename string `yaml:"filename" json:"filename" default:"run.log"`
}

// SystemConfig config of baetyl system
type SystemConfig struct {
	Certificate Certificate  `yaml:"cert,omitempty" json:"cert,omitempty" `
	Function    ClientConfig `yaml:"function,omitempty" json:"function,omitempty"`
	Core        ClientConfig `yaml:"core,omitempty" json:"core,omitempty"`
	Broker      MqttConfig   `yaml:"broker,omitempty" json:"broker,omitempty"`
	Logger      log.Config   `yaml:"logger,omitempty" json:"logger,omitempty"`
}

type MinioConfig struct {
	EndPoint  string `yaml:"endPoint" json:"endPoint"`
	AccessKey string `yaml:"accessKey" json:"accessKey"`
	SecretKey string `yaml:"secretKey" json:"secretKey"`
}

type LiveHook struct {
	API       string `yaml:"api" json:"api"`
	RTMP      string `yaml:"rtmp" json:"rtmp"`
	HLS       string `yaml:"hls" json:"hls"`
	RtmpSlave string `yaml:"rtmp_slave" json:"rtmp_slave"`
	HlsSlave  string `yaml:"hls_slave" json:"hls_slave"`
}
