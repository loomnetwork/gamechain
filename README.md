# Zombie Battleground

## Build Zombie Battleground Contract

```
make deps
make
```

## Oracle

Generate a set of private and public keys by entering the following command:

```
loom genkey -k oracle-priv.key -a oracle-pub.key > oracle-key.txt
```

List the content of `oracle-key.txt` with `cat oracle-key.txt`. You should see something like this printed out to the console:

```bash
local address: 0x97A3F939B6d14fD9C0E037963d18Bb37A9B9c646
local address base64: l6P5ObbRT9nA4DeWPRi7N6m5xkY=
```

Copy base64 address and paste it into the "oracle" section of the `genesis.json` file:

```
"oracle": {
          "chainId": "default",
          "local": <PASTE_HERE_THE_ORACLE_KEY>
```

Note that, on a chain that's already running, you can update the address of the oracle with the following command:

```bash
./bin/zb-cli update_oracle default:NEW_ORACLE_ADDRESS default:CURRENT_ORACLE_ADDRESS -k oracle.priv
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

Create a keypair by entering:

```
loom genkey -k <username>-priv.key -a <username>-pub.key
# Note that setAccount and getAccount supports all fields defined in `UpsertAccountRequest`. To make example simple,
# only two fields has been used.
```

> Don't forget to replace <username> with your username.

Create the account:

```
./bin/zb-cli create_account -k <username>-priv.key -u <username> -v v25 -d "{\"image\":\"Image\", \"game_membership_tier\": 1}"
```

Verify if the account was created:

```
./bin/zb-cli get_account -k <username>-priv.key -u <username>
```

Other useful commands (not required for spinning up the game):

```
# update account transaction
./bin/zb-cli update_account -k <username>-priv.key -u <username> -v "{\"image\":\"Image2\", \"game_membership_tier\": 2}"

# Get Decks
./bin/zb-cli get_decks -k <username>-priv.key -u <username> -v v1

# Get Deck by id
./bin/zb-cli get_deck -k <username>-priv.key -u <username> --deckId 0

# Add Deck
./bin/zb-cli create_deck -k <username>-priv.key -u <username> -v v2 -d "{\"overlordId\":\"1\", \"name\": \"NewDeck\", \"cards\": [ {\"card_name\": \"Banshee\", \"amount\": 2}, {\"card_name\": \"Breezee\", \"amount\": 1} ]}"

# Delete Deck by id
./bin/zb-cli delete_deck -k <username>-priv.key -u <username> --deckId 0 -v v1
```

# Initial setup

Contract must be initialized with correct data to work properly.

1. Set the PlasmaChain block number from which the oracle will fetch the events:

```bash
./bin/zb-cli set_last_plasma_block_number -n 3066893 # for staging plasmachain, update accordingly otherwise
```

2. Set contract configuration

```bash
./bin/zb-cli -k oracle-priv.key contract_configuration set_fiat_purchase_contract_version -v 3 # update if contract version changes
./bin/zb-cli -k oracle-priv.key contract_configuration set_initial_fiat_purchase_txid -v 85070591730234615865843651858142052964 # for dev, shift the start accordingly to the last used txid
./bin/zb-cli -k oracle-priv.key contract_configuration set_card_collection_sync_data_version -v v25 # data version to use for card sync, update accordingly to the current data version
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