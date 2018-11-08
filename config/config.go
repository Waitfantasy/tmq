package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var configFileName string

// redis 对应配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// nsq 对应配置
type NsqConfig struct {
	Addr  string `yaml:"addr"`
	Async bool   `yaml:"async"`
	Delay int    `yaml:"delay"`
}

// grpc 对应配置
type GRpcConfig struct {
	Addr       string `yaml:"addr"`
	EnableTLS  bool   `yaml:"enable_tls"`
	CertFile   string `yaml:"cert_file"`
	KeyFile    string `yaml:"key_file"`
	ServerName string `yaml:"server_name"`
}

// http 对应配置
type HttpConfig struct {
	Addr       string `yaml:"addr"`
	EnableTLS  bool   `yaml:"enableTls"`
	CaFile     string `yaml:"caFile"`
	CertFile   string `yaml:"certFile"`
	KeyFile    string `yaml:"keyFile"`
	ClientAuth bool   `yaml:"clientAuth"`
}

// unicorn id 服务对应配置
type IdServiceConfig struct {
	Addr       string `yaml:"addr"`
	EnableTLS  bool   `yaml:"enable_tls"`
	CertFile   string `yaml:"cert_file"`
	KeyFile    string `yaml:"key_file"`
	ServerName string `yaml:"server_name"`
}

type LoggerConfig struct {
	Level      string `yaml:"level"`
	Output     string `yaml:"output"`
	Split      bool   `yaml:"split"`
	FilePath   string `yaml:"file_path"`
	FilePrefix string `yaml:"file_prefix"`
	FileSuffix string `yaml:"file_suffix"`
}

type Config struct {
	Mq         string           `yaml:"mq"`
	Persistent string           `yaml:"persistent"`
	MysqlDSN   string           `yaml:"mysql_dsn"`
	Logger     *LoggerConfig    `yaml:"log"`
	GRpc       *GRpcConfig      `yaml:"grpc"`
	Http       *HttpConfig      `yaml:"http"`
	Redis      *RedisConfig     `yaml:"redis"`
	Nsq        *NsqConfig       `yaml:"nsq"`
	IdService  *IdServiceConfig `yaml:"id_service"`
}

func ParseConfigData(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func ParseConfigFile(fileName string) (*Config, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	configFileName = fileName

	return ParseConfigData(data)
}
