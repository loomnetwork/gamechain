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

DOC_IMAGE=gcr.io/robotic-catwalk-188706/gamechain-logger

chmod +x bin/gamechain-logger

echo "Building $DOC_IMAGE"
docker build -t $DOC_IMAGE:latest -f Dockerfile .
docker tag $DOC_IMAGE:$BUILD_NUMBER $DOC_IMAGE:latest

echo "Pushing $DOC_IMAGE to google container registry"
gcloud docker -- push $DOC_IMAGE:$BUILD_NUMBER
gcloud docker -- push $DOC_IMAGE:latest

# Docker image for gamechain-oracle

DOC_IMAGE_ORACLE=gcr.io/robotic-catwalk-188706/gamechain-oracle

chmod +x bin/gcoracle

echo "Building $DOC_IMAGE_ORACLE"
docker build -t $DOC_IMAGE_ORACLE:latest -f Dockerfile_gcoracle .
docker tag $DOC_IMAGE:$BUILD_NUMBER $DOC_IMAGE:latest

echo "Pushing $DOC_IMAGE_ORACLE to google container registry"
gcloud docker -- push $DOC_IMAGE_ORACLE:$BUILD_NUMBER
gcloud docker -- push $DOC_IMAGE_ORACLE:latest
