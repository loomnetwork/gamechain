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
# Note that setAccount and getAccount supports all fields defined in `UpsertAccountRequest`. To make example simple,
# only two fields has been used.

cd /path/to/project/directory

# create a key pair
LOOM_CMDPLUGINDIR=cmds/ loom genkey -k priv

# send a create account tx
LOOM_CMDPLUGINDIR=cmds/ ./bin/zb-cli createAccount -k priv -u loom -v "{\"image\":\"Image\", \"game_membership_tier\": 1}"

# get account static call
LOOM_CMDPLUGINDIR=cmds/ ./bin/zb-cli getAccount -k priv -u loom

#update account transaction
LOOM_CMDPLUGINDIR=cmds/ ./bin/zb-cli setAccount -k priv -u loom -v "{\"image\":\"Image2\", \"game_membership_tier\": 2}"
```
