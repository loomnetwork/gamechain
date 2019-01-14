#!/bin/bash

set -ex

REV=`git rev-parse --short HEAD`

export RELEASE=$REV

DOC_IMAGE=gcr.io/robotic-catwalk-188706/gamechain-logger:$REV

export GOPATH=`pwd`

mkdir -p $GOPATH/bin
export PATH=$PATH:$GOPATH/bin

go get github.com/loomnetwork/go-loom

cd ${GOPATH}/src/github.com/loomnetwork/gamechain
make deps
make
make gamechain-logger
make test

chmod +x bin/gamechain-logger

echo "sending $DOC_IMAGE"
docker build -t $DOC_IMAGE -f Dockerfile .

echo "pushing to google container registry"
gcloud docker -- push  $DOC_IMAGE

ENV=development # TODO remove this
echo "sed on k8s/${ENV}/deployment.yaml"
sed -i 's/%REV%/'"$REV"'/g' k8s/${ENV}/deployment.yaml

echo "kube apply deployment"
kubectl apply -f k8s/${ENV}/deployment.yaml  --kubeconfig=/var/lib/jenkins/${ENV}_kube_config.yaml
