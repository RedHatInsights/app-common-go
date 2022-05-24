package v1

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func NewDefaultKafkaProducer() (*kafka.Producer, error) {
	cfg, err := kafkaGenerateConfig(LoadedConfig.Kafka)
	if err != nil {
		return nil, err
	}
	return kafka.NewProducer(cfg)
}

func NewKafkaProducer(cfg *kafka.ConfigMap) (*kafka.Producer, error) {
	return kafka.NewProducer(cfg)
}

func setSASLKafkaConfig(cfgSource *KafkaConfig, cfgDest *kafka.ConfigMap) {
	(*cfgDest)["security.protocol"] = *cfgSource.Brokers[0].Sasl.SecurityProtocol
	(*cfgDest)["sasl.mechanism"] = *cfgSource.Brokers[0].Sasl.SaslMechanism
	(*cfgDest)["sasl.username"] = *cfgSource.Brokers[0].Sasl.Username
	(*cfgDest)["sasl.password"] = *cfgSource.Brokers[0].Sasl.Password
}

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
