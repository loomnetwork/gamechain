LOOMCHAIN_DIR=${GOPATH}/src/github.com/loomnetwork/loomchain

set -e
set +x

(cd ..; make oracle-abigen)
(cd ..; make bin/gcoracle)

# create oracle-plasmachain-private.key and oracle-plasmachain-public.key with "loom genkey"
${LOOMCHAIN_DIR}/gamechain genkey --private_key ../bin/oracle-plasmachain-private.key --public_key ../bin/oracle-plasmachain-public.key

../bin/gcoracle \
--oracle-log-destination file://- \
--oracle-startup-delay 0 \
--plasmachain-poll-interval 1 \
--plasmachain-private-key `cat ../bin/oracle-plasmachain-private.key` \
--gamechain-private-key `cat ../oracle.priv` \
--plasmachain-chain-id default \
--plasmachain-event-uri ws://test-z-us1.dappchains.com/queryws \
--plasmachain-read-uri http://test-z-us1.dappchains.com:80/query \
--plasmachain-write-uri http://test-z-us1.dappchains.com:80/rpc \
--plasmachain-zbgcard-contract-hex-address 0x2658d8c94062227d17a4ba61adb166e152369de3