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
export ZB_CLI="/path/to/zb-cli"
export LOOM_BIN="/path/to/loom/binary"

# Note that setAccount and getAccount supports all fields defined in `UpsertAccountRequest`. To make example simple,
# only two fields has been used.

# create a key pair
LOOM_CMDPLUGINDIR=cmds/ $LOOM_BIN genkey -k priv

# send a create account tx
LOOM_CMDPLUGINDIR=cmds/ $ZB_CLI createAccount -k priv -u loom -v "{\"image\":\"Image\", \"game_membership_tier\": 1}"

# get account static call
LOOM_CMDPLUGINDIR=cmds/ $ZB_CLI getAccount -k priv -u loom

#update account transaction
LOOM_CMDPLUGINDIR=cmds/ $ZB_CLI setAccount -k priv -u loom -v "{\"image\":\"Image2\", \"game_membership_tier\": 2}"
```
