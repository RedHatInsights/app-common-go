package v1

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestClientLoad(t *testing.T) {
	if LoadedConfig == nil {
		log.Fatal("Config didn't load in init()")
	}
	if len(LoadedConfig.Kafka.Brokers) < 1 {
		log.Fatal("Kafka brokers not loaded")
	}
	if *(LoadedConfig.Kafka.Brokers[0].Port) != 27015 {
		log.Fatal("Kafka port was not loaded")
	}
	val, ok := KafkaTopics["originalName"]
	if !ok {
		log.Fatal("Kafka Topic not found")
	}
	if val.Name != "someTopic" {
		log.Fatal("Wrong topic name")
	}
	bucket, ok := ObjectBuckets["reqname"]
	if !ok {
		log.Fatal("Object bucket not found")
	}
	if bucket.Name != "name" {
		log.Fatal("Wrong bucket name")
	}
	if KafkaServers[0] != "broker-host:27015" {
		log.Fatal("Wrong broker host")
	}
	if IsClowderEnabled() == false {
		log.Fatal("Should be true if env var ACG_CONFIG is present")
	}
	if LoadedConfig.FeatureFlags.Hostname != "ff-server.server.example.com" {
		log.Fatal("Wrong feature flag hostname")
	}
	if DependencyEndpoints["app1"]["endpoint1"].Port != 8000 {
		log.Fatal("endpoint had wrong port")
	}
	if DependencyEndpoints["app2"]["endpoint2"].Name != "endpoint2" {
		log.Fatal("endpoint had wrong name")
	}
	if PrivateDependencyEndpoints["app1"]["endpoint1"].Port != 10000 {
		log.Fatal("endpoint had wrong port")
	}
	if PrivateDependencyEndpoints["app2"]["endpoint2"].Name != "endpoint2" {
		log.Fatal("endpoint had wrong name")
	}

	rdsFilename, err := LoadedConfig.RdsCa()

	if err != nil {
		log.Fatal("error in creating RDSCa file")
	}

	content, err := ioutil.ReadFile(rdsFilename)
	if err != nil {
		log.Fatal("error reading ca")
	}

	if string(content) != *LoadedConfig.Database.RdsCa {
		log.Fatal("ca didn't match")
	}
}
