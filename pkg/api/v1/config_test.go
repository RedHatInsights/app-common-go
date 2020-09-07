package v1

import (
	"log"
	"testing"
)

func TestClientLoad(t *testing.T) {
	if LoadedConfig.Kafka.Brokers[0].Port != 27015 {
		log.Fatal("Kafka port was not loaded")
	}
}
