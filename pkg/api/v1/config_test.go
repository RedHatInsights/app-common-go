package v1

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientLoad(t *testing.T) {
	assert.NotNil(t, LoadedConfig, "Config didn't load in init()")
	assert.Len(t, LoadedConfig.Kafka.Brokers, 1, "Kafka brokers not loaded")
	assert.Equal(t, 27015, *(LoadedConfig.Kafka.Brokers[0].Port), "Kafka port was not loaded")
	assert.Contains(t, KafkaTopics, "originalName", "Kafka Topic not found")
	assert.Equal(t, "someTopic", KafkaTopics["originalName"].Name, "Wrong topic name")
	assert.Contains(t, ObjectBuckets, "reqname", "ObjectBucket not found")
	assert.Equal(t, "name", ObjectBuckets["reqname"].Name, "Wrong bucket name")

	assert.ElementsMatch(t, []string{"broker-host:27015"}, KafkaServers)
	assert.True(t, IsClowderEnabled(), "Should be true if env var ACG_CONFIG is present")

	assert.Equal(t, "ff-server.server.example.com", LoadedConfig.FeatureFlags.Hostname, "Wrong feature flag hostname")
	assert.Equal(t, "http", string(LoadedConfig.FeatureFlags.Scheme), "Wrong feature flag scheme")

	assert.Equal(t, "/foo/bar", *(LoadedConfig.TlsCAPath))

	assert.Equal(t, 8000, DependencyEndpoints["app1"]["endpoint1"].Port, "endpoint had wrong port")
	assert.Equal(t, "endpoint2", DependencyEndpoints["app2"]["endpoint2"].Name, "endpoint had wrong name")
	assert.Equal(t, 10000, PrivateDependencyEndpoints["app1"]["endpoint1"].Port, "endpoint had wrong port")
	assert.Equal(t, "endpoint2", PrivateDependencyEndpoints["app2"]["endpoint2"].Name, "endpoint had wrong name")

	rdsFilename, err := LoadedConfig.RdsCa()
	assert.Nil(t, err, "error in creating RDSCa file")
	content, err := os.ReadFile(rdsFilename)
	assert.Nil(t, err, "error reading ca")
	assert.Equal(t, *LoadedConfig.Database.RdsCa, string(content), "rds ca didn't match")

	kafkaFilename, err := LoadedConfig.KafkaCa(LoadedConfig.Kafka.Brokers[0])
	assert.Nil(t, err, "error in creating KafkaCa file")
	content, err = os.ReadFile(kafkaFilename)
	assert.Nil(t, err, "error reading ca")
	assert.Equal(t, *LoadedConfig.Kafka.Brokers[0].Cacert, string(content), "kafka ca didn't match")

	kafkaFilename, err = LoadedConfig.KafkaCa()
	assert.Nil(t, err, "error in creating KafkaCa file")
	content, err = os.ReadFile(kafkaFilename)
	assert.Nil(t, err, "error reading ca")
	assert.Equal(t, *LoadedConfig.Kafka.Brokers[0].Cacert, string(content), "kafka ca didn't match")
}

func TestEmptyRDSCa(t *testing.T) {
	cfg, err := loadConfig("../../../tests/nordsca.json")
	if err != nil {
		log.Fatalf("can't load config: %s", err)
	}

	_, err = cfg.RdsCa()

	if err == nil {
		log.Fatal("error should have been created")
	}

}
