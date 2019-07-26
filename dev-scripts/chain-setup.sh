set -e

ZB_CARD_META_DATA_DIR=${GOPATH}/src/github.com/loomnetwork/zb_card_meta_data
PLASMACHAIN_LAST_BLOCK=3294410 # for staging, update to avoid scanning long ranges every time

echo "--- Updating data"

pushd ${ZB_CARD_META_DATA_DIR}
./update_localhost.sh
popd

echo "--- Data updated"

echo "--- Update configuration"

../bin/zb-cli -k ../oracle.priv contract_configuration set_fiat_purchase_contract_version -v 3
../bin/zb-cli -k ../oracle.priv contract_configuration set_initial_fiat_purchase_txid -v `cat initialFiatPurchaseTxId.txt`

# move initial tx id by 30 each time

initialFiatPurchaseTxId=`cat initialFiatPurchaseTxId.txt`
echo "${initialFiatPurchaseTxId} + 30" | bc > initialFiatPurchaseTxId.txt

echo "--- Configuration updated"

echo "--- Setting PlasmaChain last block"
../bin/zb-cli -k ../oracle.priv set_last_plasma_block_number -n ${PLASMACHAIN_LAST_BLOCK}
echo "--- PlasmaChain last block set"