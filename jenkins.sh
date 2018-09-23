#!/bin/bash

#export GOPATH=`pwd`
export GOPATH=/var/lib/jenkins/workspace/zombie_battleground_pr
export LOOM_VER=404

mkdir -p $GOPATH/bin
export PATH=$PATH:$GOPATH/bin:/var/lib/jenkins/workspace/zombie_battleground_pr/bin #TODO make these paths less hard coded

wget https://private.delegatecall.com/loom/linux/build-${LOOM_VER}/loom -O ${GOPATH}/bin/loom
go get github.com/loomnetwork/go-loom

make deps
make
make test