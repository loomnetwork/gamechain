set -e

LOOMCHAIN_DIR=${GOPATH}/src/github.com/loomnetwork/loomchain
LOOM_NAME=gamechain
export RL_DEBUG=1
export RL_PURCHASE_GATEWAY_PRIVATE_KEY=`cat purchaseGatewayPrivateKeyHex.txt`
export CONTRACT_LOG_LEVEL=debug 
export CONTRACT_LOG_DESTINATION="file://-"

if [ $# -eq 0 ]; then
  rm -rf ../app.db
  rm -rf ../evm.db
  rm -rf ../chaindata
  rm -rf ../receipts_db
fi

#make build

pushd ${LOOMCHAIN_DIR}
make gamechain
popd

(cd ..; make cli)

set +e
if [ $# -eq 0 ]; then
  (cd ..; ${LOOMCHAIN_DIR}/${LOOM_NAME} init)
fi
set -e

cp ./config.toml ../chaindata/config/config.toml

(sleep 5; ./chain-setup.sh) &
(cd ..; ${LOOMCHAIN_DIR}/${LOOM_NAME} run)