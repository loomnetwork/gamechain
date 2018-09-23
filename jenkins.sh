#!/bin/bash

export GOPATH=`pwd`
export LOOM_VER=404

mkdir -p $GOPATH/bin
export PATH=$PATH:$GOPATH/bin

wget https://private.delegatecall.com/loom/linux/build-${LOOM_VER}/loom -O ${GOPATH}/bin/loom
go get github.com/loomnetwork/go-loom
cd ${GOPATH}/src/github.com/loomnetwork/zombie_battleground

make deps
make
make test