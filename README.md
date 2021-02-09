app-common-go
=============

Simple client access library for the config for the Clowder operator.

Usage
-----

To access configuration, see the following example

```
import (
        clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

func main() {
    if clowder.IsClowderEnabled() {
        fmt.Printf("Public Port: %s", clowder.LoadedConfig.PublicPort)
    }
}
```

The ``clowder`` library also comes with several other helpers

* ``clowder.LoadedConfig.RdsCa()`` - creates a temporary file with the RDSCa and 
  returns the filename.
* ``clowder.KafkaTopics`` - returns a map of KafkaTopics using the requestedName
  as the key and the topic object as the value.
* ``clowder.KafkaServers`` - returns a list of Kafka Broker URLs.
* ``clowder.ObjectBuckets`` - returns a list of ObjectBuckets using the requestedName
  as the key and the bucket object as the value.
* ``clowder.DependencyEndpoints`` - returns a nested map using \[appName\]\[deploymentName\] 
  for the public services of requested applications. 
* ``clowder.PrivateDependencyEndpoints`` - returns a nested map using \[appName\]\[deploymentName\] 
  for the private services of requested applications.

Testing
-------

`ACG_CONFIG="../../../test.json" go test -v ./pkg/api/v1/...`
