package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/ruixiaoedu/hiot/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
)

const (
	DefaultWorkNum = 4096
	DefaultHost    = "0.0.0.0"
	DefaultPort    = "1883"
)

type Config struct {
	WorkerNum int `json:"workerNum"`

	Host string `json:"host"`
	Port string `json:"port"`

	Cluster RouteInfo `json:"cluster"`
	Router  string    `json:"router"`

	TlsHost string  `json:"tlsHost"`
	TlsPort string  `json:"tlsPort"`
	TlsInfo TLSInfo `json:"tlsInfo"`

	WsPath string `json:"wsPath"`
	WsPort string `json:"wsPort"`
	WsTLS  bool   `json:"wsTLS"`

	Debug bool `json:"debug"`

	Database Database `json:"repository"`
	Manage   Manage   `json:"manage"`
}

type Database struct {
	Type string `json:"type"` // 数据库类型
	Dsn  string `json:"dsn"`  // dsn地址
}

type Manage struct {
	Port     int    `json:"port"`     // 管理用的HTTP端口
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

type RouteInfo struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type TLSInfo struct {
	Verify   bool   `json:"verify"`
	CaFile   string `json:"caFile"`
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

// NewConfigure 配置配置文件
func NewConfigure() (*Config, error) {
	// 默认配置
	var config = &Config{
		WorkerNum: DefaultWorkNum,
		Host:      DefaultHost,
		Port:      DefaultPort,
	}

	// 从 flag 获取配置
	if err := config.LoadFlag(); err != nil {
		return nil, err
	}

	// 从文件中获取配置
	if err := config.LoadFile(); err != nil {
		return nil, err
	}

	// 检查配置是否正确
	if err := config.Check(); err != nil {
		return nil, err
	}

	return config, nil
}

// LoadFlag 解析命令行命令
func (config *Config) LoadFlag() error {

	var configFile string // 配置文件

	kingpin.Flag("config", "Config file for hmq").Short('c').
		Default("").PlaceHolder("hiot.yml").StringVar(&configFile)

	kingpin.Flag("host", "Network host to listen on").Short('h').
		Default("0.0.0.0").StringVar(&config.Host)
	kingpin.Flag("port", "Port for MQTT to listen on.").Short('p').
		Default("1883").StringVar(&config.Port)

	kingpin.Flag("ws-port", "Port for ws to listen on").StringVar(&config.WsPort)
	kingpin.Flag("ws-path", "Path for ws to listen on").StringVar(&config.WsPath)

	kingpin.Flag("cluster-port", "Cluster port from which members can connect.").StringVar(&config.Cluster.Port)
	kingpin.Flag("router", "Router who maintenance cluster info").StringVar(&config.Router)

	kingpin.Flag("worker", "Worker num to process message, perfer (client num)/10.").Short('w').
		Default("1024").IntVar(&config.WorkerNum)

	kingpin.Flag("manage-port", "Port for HTTP API to listen on.").
		Default("8080").IntVar(&config.Manage.Port)

	kingpin.Flag("debug", "Enable Debug logging.").BoolVar(&config.Debug)

	kingpin.Parse()

	if configFile != "" {
		viper.SetConfigFile(configFile)
	}

	return nil
}

// LoadFile 解析本地配置文件
func (config *Config) LoadFile() error {

	// 没有设置配置文件，加载本地指定目录
	if viper.ConfigFileUsed() == "" {
		viper.SetConfigName("hiot")
		viper.AddConfigPath("./config")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		// 如果没有配置文件，返回原始内容
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return nil
		}
		// Config file was found but another error was produced
		return err
	}

	// 配置文件存在，转化配置文件
	if err := viper.Unmarshal(config); err != nil {
		return err
	}

	return nil
}

// Check 检查配置文件是否正确
func (config *Config) Check() error {

	if config.WorkerNum == 0 {
		config.WorkerNum = 1024
	}

	if config.Port != "" {
		if config.Host == "" {
			config.Host = "0.0.0.0"
		}
	}

	if config.Cluster.Port != "" {
		if config.Cluster.Host == "" {
			config.Cluster.Host = "0.0.0.0"
		}
	}
	if config.Router != "" {
		if config.Cluster.Port == "" {
			return errors.New("cluster port is null")
		}
	}

	if config.TlsPort != "" {
		if config.TlsInfo.CertFile == "" || config.TlsInfo.KeyFile == "" {
			logger.Error("tls config error, no cert or key file.")
			return errors.New("tls config error, no cert or key file")
		}
		if config.TlsHost == "" {
			config.TlsHost = "0.0.0.0"
		}
	}

	if config.Manage.Username == "" || config.Manage.Password == "" {
		config.Manage.Username = "admin"
		config.Manage.Password = "public"
	}

	return nil
}

func NewTLSConfig(tlsInfo TLSInfo) (*tls.Config, error) {

	cert, err := tls.LoadX509KeyPair(tlsInfo.CertFile, tlsInfo.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("error parsing X509 certificate/key pair: %v", zap.Error(err))
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", zap.Error(err))
	}

	// Create TLSConfig
	// We will determine the cipher suites that we prefer.
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	// Require client certificates as needed
	if tlsInfo.Verify {
		config.ClientAuth = tls.RequireAndVerifyClientCert
	}
	// Add in CAs if applicable.
	if tlsInfo.CaFile != "" {
		rootPEM, err := ioutil.ReadFile(tlsInfo.CaFile)
		if err != nil || rootPEM == nil {
			return nil, err
		}
		pool := x509.NewCertPool()
		ok := pool.AppendCertsFromPEM([]byte(rootPEM))
		if !ok {
			return nil, fmt.Errorf("failed to parse root ca certificate")
		}
		config.ClientCAs = pool
	}

	return &config, nil
}
