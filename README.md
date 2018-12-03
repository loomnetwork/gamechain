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
./bin/zb-cli get_decks -k priv -u loom

# Get Deck by id
./bin/zb-cli get_deck -k priv -u loom --deckId 0

# Add Deck
./bin/zb-cli create_deck -k priv -u loom -v v2 -d "{\"heroId\":\"1\", \"name\": \"NewDeck\", \"cards\": [ {\"card_name\": \"Banshee\", \"amount\": 2}, {\"card_name\": \"Breezee\", \"amount\": 1} ]}"

# Delete Deck by id
./bin/zb-cli delete_deck -k priv -u loom --deckId 0
```

## Oracle

Some transactions require oracle permissions. The oracle's private key is commited in the repo. Its address is:

```
loom genkey -k oracle.priv -a oracle.pub
local address: 0x86f36D9f1BB6af96bA809d7aA7812251424641A5
local address base64: hvNtnxu2r5a6gJ16p4EiUUJGQaU=
```

The oracle will be set automatically on chain init according the genesis file. On a chain that's already running, it can be updated with:

```
./bin/zb-cli update_oracle default:NEW_ORACLE_ADDRESS default:CURRENT_ORACLE_ADDRESS -k oracle.priv
```
.
