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
