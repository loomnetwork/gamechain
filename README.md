# Zombie Battleground

## Build Zombie Battleground Contract

```
make deps
make
```

## Run with loomchain

Make sure you have [loom](github.com/loomnetwork/loomchain) binary.

Run the follwing commands in the `gamechain` directory:
```
loom init
cp zb.genesis.json genesis.json
loom run
```


## Creating account and running transactions

```
# create a key pair using loom binary
loom genkey -k priv

# Note that setAccount and getAccount supports all fields defined in `UpsertAccountRequest`. To make example simple,
# only two fields has been used.

# In gamechain directory, run:

# create account transaction
./bin/zb-cli create_account -k priv -u loom -v v2 -d "{\"image\":\"Image\", \"game_membership_tier\": 1}"

# get account static call
./bin/zb-cli get_account -k priv -u loom

# update account transaction
./bin/zb-cli update_account -k priv -u loom -v "{\"image\":\"Image2\", \"game_membership_tier\": 2}"

# Get Decks
./bin/zb-cli get_decks -k priv -u loom -v v1

# Get Deck by id
./bin/zb-cli get_deck -k priv -u loom --deckId 0

# Add Deck
./bin/zb-cli create_deck -k priv -u loom -v v2 -d "{\"overlordId\":\"1\", \"name\": \"NewDeck\", \"cards\": [ {\"card_name\": \"Banshee\", \"amount\": 2}, {\"card_name\": \"Breezee\", \"amount\": 1} ]}"

# Delete Deck by id
./bin/zb-cli delete_deck -k priv -u loom --deckId 0 -v v1
```

## Oracle

Some transactions require oracle permissions. The oracle's private key is commited in the repo. Its address is:

```bash
loom genkey -k oracle.priv -a oracle.pub
local address: 0x86f36D9f1BB6af96bA809d7aA7812251424641A5
local address base64: hvNtnxu2r5a6gJ16p4EiUUJGQaU=
```

The oracle will be set automatically on chain init according the genesis file. On a chain that's already running, it can be updated with:

```bash
./bin/zb-cli update_oracle default:NEW_ORACLE_ADDRESS default:CURRENT_ORACLE_ADDRESS -k oracle.priv
```

# Initial setup

Contract must be initialized with correct data to work properly.

1. Set the PlasmaChain block number from which the oracle will fetch the events
```bash
./bin/zb-cli set_last_plasma_block_number -n 3066893 # for staging plasmachain, update accordingly otherwise
```

2. Set contract configuration
```bash
./bin/zb-cli -k oracle.priv contract_configuration set_fiat_purchase_contract_version -v 3 # update if contract version changes
./bin/zb-cli -k oracle.priv contract_configuration set_initial_fiat_purchase_txid -v 85070591730234615865843651858142052964 # for dev, shift the start accordingly to the last used txid
./bin/zb-cli -k oracle.priv contract_configuration set_card_collection_sync_data_version -v v25 # data version to use for card sync, update accordingly to the current data version
```

## TxId Start

IAP purchases handled by Auth have TxId incrementing from 0, and can't be reused. Since Gamechain can't communicate with Auth, the solution is to split the ranges used by Auth and Gamechain for TxId.

### Production

Starts from (2^127 + 10000) = 170141183460469231731687303715884115728
Ends on (2^256 - 1)

### Staging

Starts from (2^126 + 100000000) = 85070591730234615865843651858042052864
Ends on (2^126 + 200000000)

### Development
Starts from (2^126 + 200000000) = 85070591730234615865843651858142052864
Ends on (2^126 + 300000000)

### Local Development
Starts from (2^126) = 85070591730234615865843651857942052864
Ends on (2^126 + 100000000 - 1)

3. Export the env variable which is a key to sign marketplace transaction receipts
```bash
export RL_PURCHASE_GATEWAY_PRIVATE_KEY={hex private key}
```