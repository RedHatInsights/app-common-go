#!/bin/bash

# function to check if gojsonschema command is installed and if not install it
function check_and_install_gojsonschema {
    if ! command -v go-jsonschema &> /dev/null
    then
        echo "go-jsonschema could not be found"
        echo "installing go-jsonschema"
        go install github.com/atombender/go-jsonschema@latest
    fi
}

check_and_install_gojsonschema
if [ $(command -v wget) ]; then
  wget https://raw.githubusercontent.com/RedHatInsights/clowder/master/controllers/cloud.redhat.com/config/schema.json -O pkg/api/v1/schema.json
elif [ $(command -v curl) ]; then
  curl https://raw.githubusercontent.com/RedHatInsights/clowder/master/controllers/cloud.redhat.com/config/schema.json  -o pkg/api/v1/schema.json
fi

go-jsonschema -p v1 -o pkg/api/v1/types.go pkg/api/v1/schema.json
