// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package v1

import "fmt"
import "encoding/json"

// UnmarshalJSON implements json.Unmarshaler.
func (j *CloudWatchConfig) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["accessKeyId"]; !ok || v == nil {
		return fmt.Errorf("field accessKeyId: required")
	}
	if v, ok := raw["logGroup"]; !ok || v == nil {
		return fmt.Errorf("field logGroup: required")
	}
	if v, ok := raw["region"]; !ok || v == nil {
		return fmt.Errorf("field region: required")
	}
	if v, ok := raw["secretAccessKey"]; !ok || v == nil {
		return fmt.Errorf("field secretAccessKey: required")
	}
	type Plain CloudWatchConfig
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CloudWatchConfig(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DatabaseConfig) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["adminPassword"]; !ok || v == nil {
		return fmt.Errorf("field adminPassword: required")
	}
	if v, ok := raw["adminUsername"]; !ok || v == nil {
		return fmt.Errorf("field adminUsername: required")
	}
	if v, ok := raw["hostname"]; !ok || v == nil {
		return fmt.Errorf("field hostname: required")
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name: required")
	}
	if v, ok := raw["password"]; !ok || v == nil {
		return fmt.Errorf("field password: required")
	}
	if v, ok := raw["port"]; !ok || v == nil {
		return fmt.Errorf("field port: required")
	}
	if v, ok := raw["username"]; !ok || v == nil {
		return fmt.Errorf("field username: required")
	}
	type Plain DatabaseConfig
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DatabaseConfig(plain)
	return nil
}

// clowdapp deployment configuration for cloud.redhat.com clowdapps
type AppConfig struct {
	// Database corresponds to the JSON schema field "database".
	Database *DatabaseConfig `json:"database,omitempty"`

	// Endpoints corresponds to the JSON schema field "endpoints".
	Endpoints []DependencyEndpoint `json:"endpoints,omitempty"`

	// InMemoryDb corresponds to the JSON schema field "inMemoryDb".
	InMemoryDb *InMemoryDBConfig `json:"inMemoryDb,omitempty"`

	// Kafka corresponds to the JSON schema field "kafka".
	Kafka *KafkaConfig `json:"kafka,omitempty"`

	// Logging corresponds to the JSON schema field "logging".
	Logging LoggingConfig `json:"logging"`

	// MetricsPath corresponds to the JSON schema field "metricsPath".
	MetricsPath string `json:"metricsPath"`

	// MetricsPort corresponds to the JSON schema field "metricsPort".
	MetricsPort int `json:"metricsPort"`

	// ObjectStore corresponds to the JSON schema field "objectStore".
	ObjectStore *ObjectStoreConfig `json:"objectStore,omitempty"`

	// WebPort corresponds to the JSON schema field "webPort".
	WebPort int `json:"webPort"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DependencyEndpoint) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["app"]; !ok || v == nil {
		return fmt.Errorf("field app: required")
	}
	if v, ok := raw["hostname"]; !ok || v == nil {
		return fmt.Errorf("field hostname: required")
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name: required")
	}
	if v, ok := raw["port"]; !ok || v == nil {
		return fmt.Errorf("field port: required")
	}
	type Plain DependencyEndpoint
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DependencyEndpoint(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ObjectStoreConfig) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["hostname"]; !ok || v == nil {
		return fmt.Errorf("field hostname: required")
	}
	if v, ok := raw["port"]; !ok || v == nil {
		return fmt.Errorf("field port: required")
	}
	type Plain ObjectStoreConfig
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ObjectStoreConfig(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *InMemoryDBConfig) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["hostname"]; !ok || v == nil {
		return fmt.Errorf("field hostname: required")
	}
	if v, ok := raw["port"]; !ok || v == nil {
		return fmt.Errorf("field port: required")
	}
	type Plain InMemoryDBConfig
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = InMemoryDBConfig(plain)
	return nil
}

// broker configuration
type BrokerConfig struct {
	// Hostname corresponds to the JSON schema field "hostname".
	Hostname string `json:"hostname"`

	// Port corresponds to the JSON schema field "port".
	Port *int `json:"port,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *BrokerConfig) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["hostname"]; !ok || v == nil {
		return fmt.Errorf("field hostname: required")
	}
	type Plain BrokerConfig
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = BrokerConfig(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ObjectStoreBucket) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name: required")
	}
	if v, ok := raw["requestedName"]; !ok || v == nil {
		return fmt.Errorf("field requestedName: required")
	}
	type Plain ObjectStoreBucket
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ObjectStoreBucket(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *TopicConfig) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name: required")
	}
	if v, ok := raw["requestedName"]; !ok || v == nil {
		return fmt.Errorf("field requestedName: required")
	}
	type Plain TopicConfig
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = TopicConfig(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *LoggingConfig) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["type"]; !ok || v == nil {
		return fmt.Errorf("field type: required")
	}
	type Plain LoggingConfig
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = LoggingConfig(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *KafkaConfig) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["brokers"]; !ok || v == nil {
		return fmt.Errorf("field brokers: required")
	}
	if v, ok := raw["topics"]; !ok || v == nil {
		return fmt.Errorf("field topics: required")
	}
	type Plain KafkaConfig
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = KafkaConfig(plain)
	return nil
}

// Cloud Watch configuration
type CloudWatchConfig struct {
	// AccessKeyId corresponds to the JSON schema field "accessKeyId".
	AccessKeyId string `json:"accessKeyId"`

	// LogGroup corresponds to the JSON schema field "logGroup".
	LogGroup string `json:"logGroup"`

	// Region corresponds to the JSON schema field "region".
	Region string `json:"region"`

	// SecretAccessKey corresponds to the JSON schema field "secretAccessKey".
	SecretAccessKey string `json:"secretAccessKey"`
}

// database configuration
type DatabaseConfig struct {
	// AdminPassword corresponds to the JSON schema field "adminPassword".
	AdminPassword string `json:"adminPassword"`

	// AdminUsername corresponds to the JSON schema field "adminUsername".
	AdminUsername string `json:"adminUsername"`

	// Hostname corresponds to the JSON schema field "hostname".
	Hostname string `json:"hostname"`

	// Name corresponds to the JSON schema field "name".
	Name string `json:"name"`

	// Password corresponds to the JSON schema field "password".
	Password string `json:"password"`

	// Port corresponds to the JSON schema field "port".
	Port int `json:"port"`

	// Username corresponds to the JSON schema field "username".
	Username string `json:"username"`
}

// Dependent service connection info
type DependencyEndpoint struct {
	// App corresponds to the JSON schema field "app".
	App string `json:"app"`

	// Hostname corresponds to the JSON schema field "hostname".
	Hostname string `json:"hostname"`

	// Name corresponds to the JSON schema field "name".
	Name string `json:"name"`

	// Port corresponds to the JSON schema field "port".
	Port int `json:"port"`
}

// In Memory DB configuration
type InMemoryDBConfig struct {
	// Hostname corresponds to the JSON schema field "hostname".
	Hostname string `json:"hostname"`

	// Password corresponds to the JSON schema field "password".
	Password *string `json:"password,omitempty"`

	// Port corresponds to the JSON schema field "port".
	Port int `json:"port"`

	// Username corresponds to the JSON schema field "username".
	Username *string `json:"username,omitempty"`
}

// kafka configuration
type KafkaConfig struct {
	// Brokers corresponds to the JSON schema field "brokers".
	Brokers []BrokerConfig `json:"brokers"`

	// Topics corresponds to the JSON schema field "topics".
	Topics []TopicConfig `json:"topics"`
}

// Logging Configuration
type LoggingConfig struct {
	// Cloudwatch corresponds to the JSON schema field "cloudwatch".
	Cloudwatch *CloudWatchConfig `json:"cloudwatch,omitempty"`

	// Type corresponds to the JSON schema field "type".
	Type string `json:"type"`
}

// object storage bucket
type ObjectStoreBucket struct {
	// AccessKey corresponds to the JSON schema field "accessKey".
	AccessKey *string `json:"accessKey,omitempty"`

	// Name corresponds to the JSON schema field "name".
	Name string `json:"name"`

	// RequestedName corresponds to the JSON schema field "requestedName".
	RequestedName string `json:"requestedName"`

	// SecretKey corresponds to the JSON schema field "secretKey".
	SecretKey *string `json:"secretKey,omitempty"`
}

// object storage configuration
type ObjectStoreConfig struct {
	// AccessKey corresponds to the JSON schema field "accessKey".
	AccessKey *string `json:"accessKey,omitempty"`

	// Buckets corresponds to the JSON schema field "buckets".
	Buckets []ObjectStoreBucket `json:"buckets,omitempty"`

	// Hostname corresponds to the JSON schema field "hostname".
	Hostname string `json:"hostname"`

	// Port corresponds to the JSON schema field "port".
	Port int `json:"port"`

	// SecretKey corresponds to the JSON schema field "secretKey".
	SecretKey *string `json:"secretKey,omitempty"`
}

// topic configuration
type TopicConfig struct {
	// ConsumerGroup corresponds to the JSON schema field "consumerGroup".
	ConsumerGroup *string `json:"consumerGroup,omitempty"`

	// Name corresponds to the JSON schema field "name".
	Name string `json:"name"`

	// RequestedName corresponds to the JSON schema field "requestedName".
	RequestedName string `json:"requestedName"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AppConfig) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["logging"]; !ok || v == nil {
		return fmt.Errorf("field logging: required")
	}
	if v, ok := raw["metricsPath"]; !ok || v == nil {
		return fmt.Errorf("field metricsPath: required")
	}
	if v, ok := raw["metricsPort"]; !ok || v == nil {
		return fmt.Errorf("field metricsPort: required")
	}
	if v, ok := raw["webPort"]; !ok || v == nil {
		return fmt.Errorf("field webPort: required")
	}
	type Plain AppConfig
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = AppConfig(plain)
	return nil
}
