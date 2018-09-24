#!/bin/bash

export GOPATH=/var/lib/jenkins/gopath-branches/gamechain-${GIT_LOCAL_BRANCH}
export LOOM_VER=404

mkdir -p ${GOPATH}/bin ; true
mkdir -p ${GOPATH}/src/github.com/loomnetwork ; true

ln -sfn `pwd` ${GOPATH}/src/github.com/loomnetwork/zombie_battleground

mkdir -p $GOPATH/bin
export PATH=$PATH:$GOPATH/bin

wget https://private.delegatecall.com/loom/linux/build-${LOOM_VER}/loom -O ${GOPATH}/bin/loom
chmod +x  ${GOPATH}/bin/loom

make deps
make
make test