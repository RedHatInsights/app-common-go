package v1

import (
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
}
