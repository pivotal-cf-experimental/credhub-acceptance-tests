#!/bin/bash

set -ex
export GOPATH=$PWD/go
export PATH=$PATH:$GOPATH/bin

cd go/src/github.com/pivotal-cf/credhub-cli
make dependencies
cd ../cred-hub-acceptance-tests
cat > config.json <<EOF
{
  "api_url": "$API_URL"
}
EOF

ginkgo -r integration
