package v1

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Generates a Kafka producer with a default config.
func NewDefaultKafkaProducer() (*kafka.Producer, error) {
	cfg, err := kafkaGenerateConfig(LoadedConfig.Kafka)
	if err != nil {
		return nil, err
	}
	return kafka.NewProducer(cfg)
}

// Generates a Kafka producer with a given config
func NewKafkaProducer(cfg *kafka.ConfigMap) (*kafka.Producer, error) {
	return kafka.NewProducer(cfg)
}

func NewDefaultKafkaConsumer(groupid string) (*kafka.Consumer, error) {
	cfg, err := kafkaGenerateConfig(LoadedConfig.Kafka)
	if err != nil {
		return nil, err
	}
	(*cfg)["group.id"] = groupid
	return kafka.NewConsumer(cfg)
}

func NewKafkaConsumer(cfg *kafka.ConfigMap) (*kafka.Consumer, error) {
	return kafka.NewConsumer(cfg)
}

func setSASLKafkaConfig(cfgSource *KafkaConfig, cfgDest *kafka.ConfigMap) {
	(*cfgDest)["security.protocol"] = *cfgSource.Brokers[0].Sasl.SecurityProtocol
	(*cfgDest)["sasl.mechanism"] = *cfgSource.Brokers[0].Sasl.SaslMechanism
	(*cfgDest)["sasl.username"] = *cfgSource.Brokers[0].Sasl.Username
	(*cfgDest)["sasl.password"] = *cfgSource.Brokers[0].Sasl.Password
}

// kafkaGenerateConfig generates a kafka.ConfigMap from a KafkaConfig struct
// used in the creation of kafka.Producer and kafka.Consumer
func kafkaGenerateConfig(cfg *KafkaConfig) (*kafka.ConfigMap, error) {
	brokers := []string{}
	for _, broker := range cfg.Brokers {
		hostPlusPort := generateBrokerString(broker.Hostname, *broker.Port)
		brokers = append(brokers, hostPlusPort)
	}
	basemap := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(brokers, ","),
	}

	if cfg.Brokers[0].Authtype != nil && (*cfg.Brokers[0].Authtype == BrokerConfigAuthtypeSasl) {
		setSASLKafkaConfig(cfg, basemap)
	}

	if cfg.Brokers[0].Cacert != nil {
		err := setCAKafkaConfig(cfg, basemap)
		if err != nil {
			return nil, err
		}
	}
	return basemap, nil
}

// setCAKafkaConfig sets the CA config for the kafka.ConfigMap
func setCAKafkaConfig(cfgSource *KafkaConfig, cfgDest *kafka.ConfigMap) error {
	filename, err := KafkaCa(cfgSource.Brokers[0])

	if err != nil {
		return err
	}
	(*cfgDest)["ssl.ca.location"] = filename

	return nil
}

func generateBrokerString(hostname string, port int) string {
	portString := strconv.Itoa(port)
	return fmt.Sprintf("%s:%s", hostname, portString)
}
