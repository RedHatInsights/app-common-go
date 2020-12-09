package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ConfigOption func(*AppConfig)

var LoadedConfig *AppConfig
var KafkaTopics map[string]TopicConfig
var DependencyEndpoints map[string]map[string]DependencyEndpoint
var ObjectBuckets map[string]ObjectStoreBucket
var KafkaServers []string

func loadConfig(filename string) (*AppConfig, error) {
	var appConfig AppConfig
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &appConfig)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	return &appConfig, nil
}

func init() {
	loadedConfig, err := loadConfig(os.Getenv("ACG_CONFIG"))
	if err != nil {
		fmt.Println(err)
		return
	}
	LoadedConfig = loadedConfig
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
	if LoadedConfig.Kafka != nil {
		for _, broker := range LoadedConfig.Kafka.Brokers {
			KafkaServers = append(KafkaServers, fmt.Sprintf("%s:%d", broker.Hostname, *broker.Port))
		}
	}
}
