#!/bin/bash

set -ex

cd ${WORKSPACE}/src/github.com/loomnetwork/gamechain

if [ "${BRANCH}" = "master" ]; then
	gsutil cp bin/zb-cli gs://private.delegatecall.com/zombie_battleground/linux/latest/zb-cli
fi
gsutil cp bin/zb-cli gs://private.delegatecall.com/zombie_battleground/linux/build-${BUILD_NUMBER}/zb-cli

if [ -f "contracts/zombiebattleground.1.0.0" ]; then
	if [ "${BRANCH}" = "master" ]; then
		gsutil cp contracts/zombiebattleground.1.0.0 gs://private.delegatecall.com/zombie_battleground/linux/latest/zombiebattleground.1.0.0
	fi
	gsutil cp contracts/zombiebattleground.1.0.0 gs://private.delegatecall.com/zombie_battleground/linux/build-${BUILD_NUMBER}/zombiebattleground.1.0.0
fi

if [ -f "contracts/zombiebattleground.so.1.0.0" ]; then
	if [ "${BRANCH}" = "master" ]; then
		gsutil cp contracts/zombiebattleground.so.1.0.0 gs://private.delegatecall.com/zombie_battleground/linux/latest/zombiebattleground.so.1.0.0
	fi
	gsutil cp contracts/zombiebattleground.so.1.0.0 gs://private.delegatecall.com/zombie_battleground/linux/build-${BUILD_NUMBER}/zombiebattleground.so.1.0.0
fi

# Since the repository is not public, we need the genesis copied to storage so that we can deploy
if [ "${BRANCH}" = "master" ]; then
	gsutil cp zb.genesis.json gs://private.delegatecall.com/zombie_battleground/linux/latest/genesis.json
fi
gsutil cp zb.genesis.json gs://private.delegatecall.com/zombie_battleground/linux/build-${BUILD_NUMBER}/genesis.json

if [ "${BRANCH}" = "master" ]; then
	gsutil cp zb.genesis.json gs://private.delegatecall.com/zombie_battleground/linux/latest/genesis.json.v1
fi
gsutil cp zb.genesis.json gs://private.delegatecall.com/zombie_battleground/linux/build-${BUILD_NUMBER}/genesis.json.v1

if [ "${BRANCH}" = "master" ]; then
	gsutil cp zb.v2.genesis.json gs://private.delegatecall.com/zombie_battleground/linux/latest/genesis.json.v2
fi
gsutil cp zb.v2.genesis.json gs://private.delegatecall.com/zombie_battleground/linux/build-${BUILD_NUMBER}/genesis.json.v2

# Custom loom
if [ "${BRANCH}" = "master" ]; then
	gsutil cp ${WORKSPACE}/bin/loom gs://private.delegatecall.com/zombie_battleground/linux/latest/loom
fi
gsutil cp ${WORKSPACE}/bin/loom gs://private.delegatecall.com/zombie_battleground/linux/build-${BUILD_NUMBER}/loom

# Custom gamechain
if [ "${BRANCH}" = "master" ]; then
	gsutil cp ${WORKSPACE}/bin/gamechain gs://private.delegatecall.com/zombie_battleground/linux/latest/gamechain
fi
gsutil cp ${WORKSPACE}/bin/gamechain gs://private.delegatecall.com/zombie_battleground/linux/build-${BUILD_NUMBER}/gamechain
