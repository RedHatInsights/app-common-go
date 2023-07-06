#!/bin/bash

# function to check if gojsonschema command is installed and if not install it
function check_and_install_gojsonschema {
    if ! command -v gojsonschema &> /dev/null
    then
        echo "gojsonschema could not be found"
        echo "installing gojsonschema"
        go install github.com/atombender/go-jsonschema/cmd/gojsonschema@latest
    fi
}

check_and_install_gojsonschema
cd pkg/api/v1
wget https://raw.githubusercontent.com/RedHatInsights/clowder/master/controllers/cloud.redhat.com/config/schema.json -O schema.json
gojsonschema -p v1 -o types.go schema.json
