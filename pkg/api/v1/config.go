package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type LoggingConfig struct {
	Type       string           `json:"type"`
	CloudWatch CloudWatchConfig `json:"cloudwatch,omitempty"`
}

type CloudWatchConfig struct {
	AccessKeyID     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	Region          string `json:"region"`
	LogGroup        string `json:"logGroup"`
}

type BrokerConfig struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

type TopicConfig struct {
	Name              string `json:"name"`
	ConsumerGroupName string `json:"consumerGroup,omitempty"`
}

type KafkaConfig struct {
	Brokers []BrokerConfig `json:"brokers"`
	Topics  []TopicConfig  `json:"topics"`
}

type ObjectStoreConfig struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Endpoint  string `json:"endpoint"`
}

type DatabaseConfig struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
	Port     int32  `json:"port"`
}

type AppConfig struct {
	WebPort     int32             `json:"webPort"`
	MetricsPort int32             `json:"metricsPort"`
	MetricsPath string            `json:"metricsPath"`
	Logging     LoggingConfig     `json:"logging"`
	Kafka       KafkaConfig       `json:"kafka"`
	Database    DatabaseConfig    `json:"database"`
	ObjectStore ObjectStoreConfig `json:"objectStore"`
}

type ConfigOption func(*AppConfig)

var LoadedConfig *AppConfig

func loadConfig(filename string) *AppConfig {
	var appConfig AppConfig
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &appConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	return &appConfig
}

func init() {
	LoadedConfig = loadConfig(os.Getenv("ACG_CONFIG"))
}

func ObjectStore(store ObjectStoreConfig) ConfigOption {
	return func(c *AppConfig) {
		c.ObjectStore = store
	}
}

func Logging(logging LoggingConfig) ConfigOption {
	return func(c *AppConfig) {
		c.Logging = logging
	}
}

func Kafka(kc KafkaConfig) ConfigOption {
	return func(c *AppConfig) {
		c.Kafka = kc
	}
}

func Database(dc DatabaseConfig) ConfigOption {
	return func(c *AppConfig) {
		c.Database = dc
	}
}

func Web(port int32) ConfigOption {
	return func(c *AppConfig) {
		c.WebPort = port
	}
}

func Metrics(path string, port int32) ConfigOption {
	return func(c *AppConfig) {
		c.MetricsPath = path
		c.MetricsPort = port
	}
}

func New(opts ...ConfigOption) *AppConfig {
	c := &AppConfig{
		WebPort:     int32(8080),
		MetricsPort: int32(9090),
		MetricsPath: "/metrics",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}
