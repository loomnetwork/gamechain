set -e

LOOMCHAIN_DIR=${GOPATH}/src/github.com/loomnetwork/loomchain

if [ $# -eq 0 ]; then
  rm -rf ./app.db
  rm -rf ./chaindata
fi

make
make cli

set +e
${LOOMCHAIN_DIR}/loom init
set -e

cp ./config.toml ./chaindata/config/config.toml

${LOOMCHAIN_DIR}/loom run