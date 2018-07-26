# Zombie Battleground

## Build Zombie Battleground Contract

```
make deps
make
```


## Run with loomchain

Make sure you have [loom](github.com/loomnetwork/loomchain) binary.

Run the follwing commands in the `zombie_battleground` directory:
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

# In zombie_battleground directory, run:

# create account transaction
./bin/zb-cli create_account -k priv -u loom -v "{\"image\":\"Image\", \"game_membership_tier\": 1}"

# get account static call
./bin/zb-cli get_account -k priv -u loom

# update account transaction
./bin/zb-cli update_account -k priv -u loom -v "{\"image\":\"Image2\", \"game_membership_tier\": 2}"

# Get Decks
./bin/zb-cli get_decks -k priv -u loom

# Get Deck by name
./bin/zb-cli get_deck -k priv -u loom -d Default

# Add Deck
./bin/zb-cli create_deck -k priv -u loom -v "{\"heroId\":\"1\", \"name\": \"NewDeck\", \"cards\": [ {\"card_id\": 1, \"amount\": 2}, {\"card_id\": 2, \"amount\": 1} ]}"

# Delete Deck
./bin/zb-cli delete_deck -k priv -u loom -d NewDeck
```



