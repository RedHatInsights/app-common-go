#!/bin/bash
cd pkg/api/v1
wget https://raw.githubusercontent.com/RedHatInsights/clowder/master/controllers/cloud.redhat.com/config/schema.json -O schema.json
gojsonschema -p v1 -o types.go schema.json
