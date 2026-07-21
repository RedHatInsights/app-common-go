package v1

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	assert.Equal(t, "app1-api-path", DependencyEndpoints["app1"]["endpoint1"].ApiPath, "endpoint1 had wrong port")
	assert.Equal(t, "app2-api-path", DependencyEndpoints["app2"]["endpoint2"].ApiPath, "endpoint2 had wrong port")

	assert.Equal(t, 8000, DependencyEndpoints["app1"]["endpoint1"].Port, "endpoint had wrong port")
	assert.Equal(t, "endpoint2", DependencyEndpoints["app2"]["endpoint2"].Name, "endpoint had wrong name")
	assert.Equal(t, 10000, PrivateDependencyEndpoints["app1"]["endpoint1"].Port, "endpoint had wrong port")
	assert.Equal(t, "endpoint2", PrivateDependencyEndpoints["app2"]["endpoint2"].Name, "endpoint had wrong name")

	rdsFilename, err := LoadedConfig.RdsCa()
	assert.NoError(t, err, "error in creating RDSCa file")
	content, err := os.ReadFile(rdsFilename)
	assert.NoError(t, err, "error reading ca")
	assert.Equal(t, *LoadedConfig.Database.RdsCa, string(content), "rds ca didn't match")

	kafkaFilename, err := LoadedConfig.KafkaCa(LoadedConfig.Kafka.Brokers[0])
	assert.NoError(t, err, "error in creating KafkaCa file")
	content, err = os.ReadFile(kafkaFilename)
	assert.NoError(t, err, "error reading ca")
	assert.Equal(t, *LoadedConfig.Kafka.Brokers[0].Cacert, string(content), "kafka ca didn't match")

	kafkaFilename, err = LoadedConfig.KafkaCa()
	assert.NoError(t, err, "error in creating KafkaCa file")
	content, err = os.ReadFile(kafkaFilename)
	assert.NoError(t, err, "error reading ca")
	assert.Equal(t, *LoadedConfig.Kafka.Brokers[0].Cacert, string(content), "kafka ca didn't match")

	assert.Equal(t, "testing", *LoadedConfig.Hostname, "top level hostname didn't match")
}

func TestEmptyRDSCa(t *testing.T) {
	cfg, err := LoadConfig("testdata/nordsca.json")
	require.NoErrorf(t, err, "can't load config: %s", err)

	path, err := cfg.RdsCa()
	require.Empty(t, path)
	require.Error(t, err, "error should have been created")
}

func TestV2DependencyEndpoints(t *testing.T) {
	// Verify V2 endpoints were parsed from test.json
	assert.NotNil(t, DependencyEndpointsV2, "V2 public endpoints should be populated")
	assert.NotNil(t, PrivateDependencyEndpointsV2, "V2 private endpoints should be populated")

	// Test app1 service1: in-cluster (authenticated: false, no CA cert)
	endpoint, ok := DependencyEndpointsV2["app1"]["service1"]
	assert.True(t, ok, "app1/service1 should exist in V2 endpoints")
	assert.Equal(t, "http://app1-service1.svc:8080", endpoint.Uri, "URI should match")
	assert.False(t, endpoint.Authenticated, "in-cluster endpoint should have authenticated=false")
	assert.Nil(t, endpoint.CaCertificate, "in-cluster endpoint should not have CA certificate")

	// Test app1 service2: cross-cluster (authenticated: true, with CA cert)
	endpoint, ok = DependencyEndpointsV2["app1"]["service2"]
	assert.True(t, ok, "app1/service2 should exist in V2 endpoints")
	assert.Equal(t, "https://app1-service2.example.com:8443", endpoint.Uri, "URI should match")
	assert.True(t, endpoint.Authenticated, "cross-cluster endpoint should have authenticated=true")
	assert.NotNil(t, endpoint.CaCertificate, "cross-cluster endpoint should have CA certificate")
	assert.Equal(t, "/cdapp/certs/app1-service2-ca.crt", *endpoint.CaCertificate, "CA cert path should match")

	// Test app2 service1
	endpoint, ok = DependencyEndpointsV2["app2"]["service1"]
	assert.True(t, ok, "app2/service1 should exist in V2 endpoints")
	assert.Equal(t, "http://app2-service1.svc:9000", endpoint.Uri, "URI should match")

	// Test private endpoints
	endpoint, ok = PrivateDependencyEndpointsV2["app1"]["privateService1"]
	assert.True(t, ok, "app1/privateService1 should exist in V2 private endpoints")
	assert.Equal(t, "http://app1-private.svc:10000", endpoint.Uri, "Private endpoint URI should match")
	assert.False(t, endpoint.Authenticated, "private endpoint should have authenticated=false")
}

func TestGetV2DependencyEndpoint(t *testing.T) {
	// Test successful lookup
	endpoint, ok := GetV2DependencyEndpoint("app1", "service1")
	assert.True(t, ok, "should find app1/service1")
	assert.Equal(t, "http://app1-service1.svc:8080", endpoint.Uri)

	// Test not found: non-existent app
	_, ok = GetV2DependencyEndpoint("nonexistent", "service")
	assert.False(t, ok, "should not find non-existent app")

	// Test not found: non-existent service
	_, ok = GetV2DependencyEndpoint("app1", "nonexistent")
	assert.False(t, ok, "should not find non-existent service")
}

func TestGetV2PrivateDependencyEndpoint(t *testing.T) {
	// Test successful lookup
	endpoint, ok := GetV2PrivateDependencyEndpoint("app1", "privateService1")
	assert.True(t, ok, "should find app1/privateService1")
	assert.Equal(t, "http://app1-private.svc:10000", endpoint.Uri)

	// Test not found: non-existent app
	_, ok = GetV2PrivateDependencyEndpoint("nonexistent", "service")
	assert.False(t, ok, "should not find non-existent app")

	// Test not found: non-existent service
	_, ok = GetV2PrivateDependencyEndpoint("app1", "nonexistent")
	assert.False(t, ok, "should not find non-existent service")
}

func TestMalformedV2Endpoints(t *testing.T) {
	// Load a config with malformed V2 endpoints (missing required 'authenticated' field)
	// This should not panic, but V2 endpoints should be left unpopulated due to error
	cfg, err := LoadConfig("testdata/malformed_v2.json")
	require.NoErrorf(t, err, "can't load config: %s", err)

	// Verify config loaded but without V2 endpoint data
	assert.NotNil(t, cfg, "config should load even with malformed V2 data")
	// Note: we can't directly test the globals here since they were set during package init()
	// with the valid test.json. This test demonstrates that LoadConfig itself doesn't panic.
}
