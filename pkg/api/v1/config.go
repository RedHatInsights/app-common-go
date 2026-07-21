package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type ConfigOption func(*AppConfig)

var LoadedConfig *AppConfig
var KafkaTopics map[string]TopicConfig
var DependencyEndpoints map[string]map[string]DependencyEndpoint
var PrivateDependencyEndpoints map[string]map[string]PrivateDependencyEndpoint
var DependencyEndpointsV2 map[string]map[string]DependencyEndpointV2
var PrivateDependencyEndpointsV2 map[string]map[string]DependencyEndpointV2
var ObjectBuckets map[string]ObjectStoreBucket
var KafkaServers []string

func LoadConfig(filename string) (*AppConfig, error) {
	var appConfig AppConfig
	content, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &appConfig)
	if err != nil {
		return nil, err
	}
	return &appConfig, nil
}

func IsClowderEnabled() bool {
	_, ok := os.LookupEnv("ACG_CONFIG")
	return ok
}

// buildV2Endpoints parses raw V2 endpoint maps into typed DependencyEndpointV2 structs.
// Returns error if any entry is malformed (missing required fields).
func buildV2Endpoints(raw map[string]interface{}) (map[string]map[string]DependencyEndpointV2, error) {
	result := make(map[string]map[string]DependencyEndpointV2)
	for appName, appServices := range raw {
		// appServices is map[string]interface{} due to schema limitations
		appServicesMap, ok := appServices.(map[string]interface{})
		if !ok {
			continue
		}
		result[appName] = make(map[string]DependencyEndpointV2)
		for serviceName, rawEndpoint := range appServicesMap {
			// Re-marshal and unmarshal to convert raw interface{} to typed DependencyEndpointV2
			// This triggers the generated UnmarshalJSON validation (required fields: uri, authenticated)
			jsonBytes, err := json.Marshal(rawEndpoint)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal V2 endpoint %s/%s: %w", appName, serviceName, err)
			}
			var endpoint DependencyEndpointV2
			if err := json.Unmarshal(jsonBytes, &endpoint); err != nil {
				return nil, fmt.Errorf("failed to unmarshal V2 endpoint %s/%s: %w", appName, serviceName, err)
			}
			result[appName][serviceName] = endpoint
		}
	}
	return result, nil
}

// initKafkaTopics populates the KafkaTopics global from config.
func initKafkaTopics(cfg *AppConfig) {
	KafkaTopics = make(map[string]TopicConfig)
	if cfg.Kafka != nil {
		for _, topic := range cfg.Kafka.Topics {
			KafkaTopics[topic.RequestedName] = topic
		}
	}
}

// initV1DependencyEndpoints populates the V1 dependency endpoint globals from config.
func initV1DependencyEndpoints(cfg *AppConfig) {
	DependencyEndpoints = make(map[string]map[string]DependencyEndpoint)
	if cfg.Endpoints != nil {
		for _, endpoint := range cfg.Endpoints {
			if DependencyEndpoints[endpoint.App] == nil {
				DependencyEndpoints[endpoint.App] = make(map[string]DependencyEndpoint)
			}
			DependencyEndpoints[endpoint.App][endpoint.Name] = endpoint
		}
	}

	PrivateDependencyEndpoints = make(map[string]map[string]PrivateDependencyEndpoint)
	if cfg.PrivateEndpoints != nil {
		for _, endpoint := range cfg.PrivateEndpoints {
			if PrivateDependencyEndpoints[endpoint.App] == nil {
				PrivateDependencyEndpoints[endpoint.App] = make(map[string]PrivateDependencyEndpoint)
			}
			PrivateDependencyEndpoints[endpoint.App][endpoint.Name] = endpoint
		}
	}
}

// initV2DependencyEndpoints populates the V2 dependency endpoint globals from config.
func initV2DependencyEndpoints(cfg *AppConfig) {
	if cfg.DependencyEndpoints != nil && cfg.DependencyEndpoints.V2 != nil {
		v2Endpoints, err := buildV2Endpoints(cfg.DependencyEndpoints.V2)
		if err != nil {
			fmt.Println(err)
		} else {
			DependencyEndpointsV2 = v2Endpoints
		}
	}

	if cfg.PrivateDependencyEndpoints != nil && cfg.PrivateDependencyEndpoints.V2 != nil {
		v2PrivateEndpoints, err := buildV2Endpoints(cfg.PrivateDependencyEndpoints.V2)
		if err != nil {
			fmt.Println(err)
		} else {
			PrivateDependencyEndpointsV2 = v2PrivateEndpoints
		}
	}
}

// initObjectBuckets populates the ObjectBuckets global from config.
func initObjectBuckets(cfg *AppConfig) {
	ObjectBuckets = make(map[string]ObjectStoreBucket)
	if cfg.ObjectStore != nil {
		for _, bucket := range cfg.ObjectStore.Buckets {
			ObjectBuckets[bucket.RequestedName] = bucket
		}
	}
}

// initKafkaServers populates the KafkaServers global from config.
func initKafkaServers(cfg *AppConfig) {
	if cfg.Kafka != nil {
		for _, broker := range cfg.Kafka.Brokers {
			KafkaServers = append(KafkaServers, fmt.Sprintf("%s:%d", broker.Hostname, *broker.Port))
		}
	}
}

func init() {
	if !IsClowderEnabled() {
		return
	}
	loadedConfig, err := LoadConfig(os.Getenv("ACG_CONFIG"))
	if err != nil {
		fmt.Println(err)
		return
	}
	LoadedConfig = loadedConfig
	initKafkaTopics(LoadedConfig)
	initV1DependencyEndpoints(LoadedConfig)
	initV2DependencyEndpoints(LoadedConfig)
	initObjectBuckets(LoadedConfig)
	initKafkaServers(LoadedConfig)
}

// GetV2DependencyEndpoint retrieves a V2 dependency endpoint by app and service name.
// Returns (endpoint, false) if not found or V2 endpoints are not available.
func GetV2DependencyEndpoint(app, name string) (DependencyEndpointV2, bool) {
	if DependencyEndpointsV2 == nil {
		return DependencyEndpointV2{}, false
	}
	if appEndpoints, ok := DependencyEndpointsV2[app]; ok {
		if endpoint, ok := appEndpoints[name]; ok {
			return endpoint, true
		}
	}
	return DependencyEndpointV2{}, false
}

// GetV2PrivateDependencyEndpoint retrieves a V2 private dependency endpoint by app and service name.
// Returns (endpoint, false) if not found or V2 endpoints are not available.
func GetV2PrivateDependencyEndpoint(app, name string) (DependencyEndpointV2, bool) {
	if PrivateDependencyEndpointsV2 == nil {
		return DependencyEndpointV2{}, false
	}
	if appEndpoints, ok := PrivateDependencyEndpointsV2[app]; ok {
		if endpoint, ok := appEndpoints[name]; ok {
			return endpoint, true
		}
	}
	return DependencyEndpointV2{}, false
}

// RdsCa writes the RDS CA from the JSON config to a temporary file and returns
// the path
func (a AppConfig) RdsCa() (string, error) {
	return writeContent("rdsca", "rds", a.Database.RdsCa)
}

// KafkaCa writes the Kafka CA from the JSON config to a temporary file and returns
// the path
func (a AppConfig) KafkaCa(brokers ...BrokerConfig) (string, error) {
	if len(brokers) == 0 {
		if len(LoadedConfig.Kafka.Brokers) == 0 {
			return "", errors.New("no broker availabl")
		}
		brokers = LoadedConfig.Kafka.Brokers
	}
	return writeContent("kafkaca", "kafka", brokers[0].Cacert)
}

func (a AppConfig) KafkaFirstCa() (string, error) {
	if a.Kafka == nil || len(a.Kafka.Brokers) == 0 || a.Kafka.Brokers[0].Cacert == nil {
		return "", errors.New("could not find ca for first broker")
	}
	file := a.Kafka.Brokers[0].Cacert
	return writeContent("kafkaca", "kafka", file)
}

func writeContent(dir string, file string, contentString *string) (string, error) {
	dir, err := os.MkdirTemp("", dir)
	if err != nil {
		return "", err
	}

	if contentString == nil {
		return "", errors.New("no RDS available")
	}

	content := []byte(*contentString)

	tmpFile, err := os.CreateTemp(dir, file)

	if err != nil {
		return "", err
	}

	if err := os.WriteFile(tmpFile.Name(), content, 0600); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}
