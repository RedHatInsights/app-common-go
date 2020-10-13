package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type ConfigOption func(*AppConfig)

var LoadedConfig *AppConfig
var KafkaTopics map[string]TopicConfig
var DependencyEndpoints map[string]map[string]DependencyEndpoint
var ObjectBuckets map[string]ObjectStoreBucket

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
	KafkaTopics = make(map[string]TopicConfig)
	if LoadedConfig.Kafka != nil {
		for _, topic := range LoadedConfig.Kafka.Topics {
			KafkaTopics[topic.RequestedName] = topic
		}
	}
	DependencyEndpoints = make(map[string]map[string]DependencyEndpoint)
	if LoadedConfig.Endpoints != nil {
		for _, endpoint := range LoadedConfig.Endpoints {
			if DependencyEndpoints[endpoint.App] == nil {
				DependencyEndpoints[endpoint.App] = make(map[string]DependencyEndpoint)
			}
			DependencyEndpoints[endpoint.App][endpoint.Name] = endpoint
		}
	}
	ObjectBuckets = make(map[string]ObjectStoreBucket)
	if LoadedConfig.ObjectStore != nil {
		for _, bucket := range LoadedConfig.ObjectStore.Buckets {
			ObjectBuckets[bucket.RequestedName] = bucket
		}
	}
}
