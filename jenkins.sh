#!/bin/bash

set -ex

export GOPATH=`pwd`

mkdir -p $GOPATH/bin
export PATH=$PATH:$GOPATH/bin

go get github.com/loomnetwork/go-loom

cd ${GOPATH}/src/github.com/loomnetwork/gamechain
make deps
make
pushd .
cd  ${GOPATH}/src/github.com/loomnetwork/loomchain
make gamechain-cleveldb
cp gamechain $GOPATH/bin/gamechain
popd
make gamechain-logger
make bin/gcoracle
make test

# Docker image for gamechain-logger

DOC_IMAGE_LOGGER=gcr.io/robotic-catwalk-188706/gamechain-logger

chmod +x bin/gamechain-logger

echo "Building $DOC_IMAGE_LOGGER"
docker build -t $DOC_IMAGE_LOGGER:latest -f Dockerfile .
docker tag $DOC_IMAGE_LOGGER:latest $DOC_IMAGE_LOGGER:$BUILD_NUMBER

echo "Pushing $DOC_IMAGE_LOGGER to google container registry"
gcloud docker -- push $DOC_IMAGE_LOGGER:$BUILD_NUMBER
gcloud docker -- push $DOC_IMAGE_LOGGER:latest

# Docker image for gamechain-oracle

DOC_IMAGE_ORACLE=gcr.io/robotic-catwalk-188706/gamechain-oracle

chmod +x bin/gcoracle

echo "Building $DOC_IMAGE_ORACLE"
docker build -t $DOC_IMAGE_ORACLE:latest -f Dockerfile_gcoracle .
docker tag $DOC_IMAGE_ORACLE:latest $DOC_IMAGE_ORACLE:$BUILD_NUMBER

echo "Pushing $DOC_IMAGE_ORACLE to google container registry"
gcloud docker -- push $DOC_IMAGE_ORACLE:$BUILD_NUMBER
gcloud docker -- push $DOC_IMAGE_ORACLE:latest
