package broker

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/fhmq/hmq/logger"
	"github.com/fhmq/hmq/plugins/bridge"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"strings"
)

type Config struct {
	Worker   int       `json:"workerNum"`
	HTTPPort string    `json:"httpPort"`
	Host     string    `json:"host"`
	Port     string    `json:"port"`
	Cluster  RouteInfo `json:"cluster"`
	Router   string    `json:"router"`
	TlsHost  string    `json:"tlsHost"`
	TlsPort  string    `json:"tlsPort"`
	WsPath   string    `json:"wsPath"`
	WsPort   string    `json:"wsPort"`
	WsTLS    bool      `json:"wsTLS"`
	TlsInfo  TLSInfo   `json:"tlsInfo"`
	Debug    bool      `json:"debug"`
	Plugin   Plugins   `json:"plugins"`
	Database Database  `json:"database"`
}

type Database struct {
	Type   string `json:"type"`   // 数据库类型
	Dsn    string `json:"dsn"`    // dsn地址
	Extend string `json:"extend"` // 其他配置，需自行解析
}

type Plugins struct {
	Bridge bridge.BridgeMQ
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

var DefaultConfig = &Config{
	Worker: 4096,
	Host:   "0.0.0.0",
	Port:   "1883",
}

var (
	log = logger.Prod().Named("broker")
)

// ConfigureConfig 配置配置文件
func ConfigureConfig() (*Config, error) {

	// 从 flag 获取配置
	config, err := LoadFlag(DefaultConfig)
	if err != nil {
		return nil, err
	}

	// 从文件中获取配置
	config, err = LoadConfig(config)
	if err != nil {
		return nil, err
	}

	config.Plugin.Bridge = bridge.NewBridgeMQ("")

	if config.Debug {
		log = logger.Debug().Named("broker")
	}

	// 检查配置是否正确
	if err := config.check(); err != nil {
		return nil, err
	}

	return config, nil

}

// LoadFlag 解析命令行命令
func LoadFlag(config *Config) (*Config, error) {

	configFile := *kingpin.Flag("config", "Config file for hmq").Short('c').
		Default("").PlaceHolder("hiot.yml").String()
	if configFile != "" {
		viper.SetConfigFile(configFile)
	}

	kingpin.Flag("host", "Network host to listen on").Short('h').
		Default("0.0.0.0").StringVar(&config.Host)
	kingpin.Flag("port", "Port for MQTT to listen on.").Short('p').
		Default("1883").StringVar(&config.Port)

	kingpin.Flag("ws-port", "Port for ws to listen on").StringVar(&config.WsPort)
	kingpin.Flag("ws-path", "Path for ws to listen on").StringVar(&config.WsPath)

	kingpin.Flag("cluster-port", "Cluster port from which members can connect.").StringVar(&config.Cluster.Port)
	kingpin.Flag("router", "Router who maintenance cluster info").StringVar(&config.Router)

	kingpin.Flag("worker", "Worker num to process message, perfer (client num)/10.").Short('w').
		Default("1024").IntVar(&config.Worker)

	kingpin.Flag("manage-port", "Port for HTTP API to listen on.").
		Default("8080").StringVar(&config.HTTPPort)

	kingpin.Flag("debug", "Enable Debug logging.").BoolVar(&config.Debug)

	kingpin.Parse()

	return config, nil
}

// LoadConfig 解析本地配置文件
func LoadConfig(config *Config) (*Config, error) {

	viper.SetConfigName("hiot")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		// 如果没有配置文件，返回原始内容
		if strings.Index(err.Error(), "Not Found") >= 0 {
			return config, nil
		}
		return nil, err
	}

	// 配置文件存在，转化配置文件
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

// check 检查配置文件是否正确
func (config *Config) check() error {

	if config.Worker == 0 {
		config.Worker = 1024
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
			log.Error("tls config error, no cert or key file.")
			return errors.New("tls config error, no cert or key file")
		}
		if config.TlsHost == "" {
			config.TlsHost = "0.0.0.0"
		}
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
