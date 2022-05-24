package v1

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var configs map[string]*KafkaConfig

func init() {
	configs = map[string]*KafkaConfig{}
	for _, filename := range []string{"saslkafka.json", "nonsaslkafka.json"} {
		testConfig, _ := os.ReadFile(fmt.Sprintf("../../../testdata/%s", filename))
		kafkaConfig := &KafkaConfig{}
		err := json.Unmarshal(testConfig, kafkaConfig)

		if err != nil {
			panic(fmt.Errorf("could not load JSON: %w", err))
		}

		configs[filename] = kafkaConfig
	}
}

func TestSASLConfigSet(t *testing.T) {
	kafkaConfig := configs["saslkafka.json"]
	newConfig := &kafka.ConfigMap{}
	setSASLKafkaConfig(kafkaConfig, newConfig)

	assert.Equal(t, *kafkaConfig.Brokers[0].Sasl.SecurityProtocol, (*newConfig)["security.protocol"], "did not match")
	assert.Equal(t, *kafkaConfig.Brokers[0].Sasl.SaslMechanism, (*newConfig)["sasl.mechanism"], "did not match")
	assert.Equal(t, *kafkaConfig.Brokers[0].Sasl.Username, (*newConfig)["sasl.username"], "username doesn't match")
	assert.Equal(t, *kafkaConfig.Brokers[0].Sasl.Password, (*newConfig)["sasl.password"], "password doesn't match")
}

func TestBrokerStringGeneration(t *testing.T) {
	brokerHost := "broker-host"
	brokerPort := 9092
	brokerString := generateBrokerString(brokerHost, brokerPort)

	assert.Equal(t, brokerString, "broker-host:9092")
}

func TestKafkaGenerateConfigSASL(t *testing.T) {
	kafkaConfig := configs["saslkafka.json"]
	expectedHost := generateBrokerString(kafkaConfig.Brokers[0].Hostname, *kafkaConfig.Brokers[0].Port)

	config, err := kafkaGenerateConfig(kafkaConfig)

	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, config, "config should have data")

	assert.Equal(t, expectedHost, (*config)["bootstrap.servers"], "bootstrap server can't match")
	assert.Equal(t, *kafkaConfig.Brokers[0].Sasl.SaslMechanism, (*config)["sasl.mechanism"], "SASL mechanism doesn't match")
	assert.Equal(t, *kafkaConfig.Brokers[0].Sasl.Username, (*config)["sasl.username"], "username doesn't match")
	assert.Equal(t, *kafkaConfig.Brokers[0].Sasl.Password, (*config)["sasl.password"], "password doesn't match")
}

func TestKafkaGenerateConfigNonSASL(t *testing.T) {
	kafkaConfig := configs["nonsaslkafka.json"]
	expectedHost := generateBrokerString(kafkaConfig.Brokers[0].Hostname, *kafkaConfig.Brokers[0].Port)

	config, err := kafkaGenerateConfig(kafkaConfig)

	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, config, "config should have data")

	assert.Equal(t, expectedHost, (*config)["bootstrap.servers"], "bootstrap server can't match")
	assert.Nil(t, (*config)["sasl.mechanism"], "SASL mechanism doesn't match")
	assert.Nil(t, (*config)["sasl.username"], "username doesn't match")
	assert.Nil(t, (*config)["sasl.password"], "password doesn't match")
}

func TestKafkaGenerateCA(t *testing.T) {
	kafkaConfig := configs["nonsaslkafka.json"]
	expectedCAData := kafkaConfig.Brokers[0].Cacert

	config, err := kafkaGenerateConfig(kafkaConfig)

	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, config, "config should have data")

	assert.NotEmpty(t, (*config)["ssl.ca.location"], "SSL CA location is blank")

	// Load ca data from file
	filename := (*config)["ssl.ca.location"].(string)
	data, err := os.ReadFile(filename)

	assert.Nil(t, err, "error opening file")
	fileData := string(data)

	assert.Equal(t, *expectedCAData, fileData, "CA data does not match")
}

func TestKafkaProducer(t *testing.T) {
	kafkaConfig := configs["saslkafka.json"]
	cfg, _ := kafkaGenerateConfig(kafkaConfig)
	kafka, err := NewKafkaProducer(cfg)

	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, kafka, "config should have data")
}

func TestDefaultKafkaProducer(t *testing.T) {
	kafka, err := NewDefaultKafkaProducer()

	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, kafka, "config should have data")
}
