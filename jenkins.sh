#!/bin/bash

set -ex

export GOPATH=`pwd`

mkdir -p $GOPATH/bin
export PATH=$PATH:$GOPATH/bin

go get github.com/loomnetwork/go-loom

cd ${GOPATH}/src/github.com/loomnetwork/gamechain
make deps
make
make gamechain-logger
make bin/gcoracle
make test

# Docker image for gamechain-logger

REV=`git rev-parse --short HEAD`
export RELEASE=$REV
DOC_IMAGE=gcr.io/robotic-catwalk-188706/gamechain-logger:$REV

chmod +x bin/gamechain-logger

echo "sending $DOC_IMAGE"
docker build -t $DOC_IMAGE -f Dockerfile .

echo "pushing to google container registry"
gcloud docker -- push  $DOC_IMAGE

# Docker image for gamechain-oracle
DOC_IMAGE_ORACLE=gcr.io/robotic-catwalk-188706/gamechain-oracle:$REV

chmod +x bin/gcoracle

echo "sending $DOC_IMAGE_ORACLE"
docker build -t $DOC_IMAGE_ORACLE -f Dockerfile_gcoracle .

echo "pushing to google container registry"
gcloud docker -- push  $DOC_IMAGE_ORACLE

echo "sed on k8s/${ENV}/deployment.yaml"
sed -i 's/%REV%/'"$REV"'/g' k8s/${ENV}/deployment.yaml

echo "kube apply deployment"
kubectl apply -f k8s/${ENV}/deployment.yaml  --kubeconfig=/var/lib/jenkins/${ENV}_kube_config.yaml
echo "kube apply service"
kubectl apply -f k8s/${ENV}/service.yaml  --kubeconfig=/var/lib/jenkins/${ENV}_kube_config.yaml
echo "kube apply ingress"
kubectl apply -f k8s/${ENV}/ingress.yaml  --kubeconfig=/var/lib/jenkins/${ENV}_kube_config.yaml
