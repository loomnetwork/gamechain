# Gamechain Oracle

The code is mostly based on loomchain's [Oracle](https://github.com/loomnetwork/loomchain/tree/master/gateway). The upstream is Plasmachain which runs multiple Ethereum contracts using loomchain binary. The downstream is Gamechain which is a Go contracts.

The idea of this Oracle is to fetch all the events related to open pack and sycn to Gamechain. Ideally, we should have all of our cards we open on [loom.games](https://loom.games/) to Gamehchain to actually play the Zombie Battleground.

## ABI

To fetch data from Plasmachain we need contract's abi to generate ethcontracts. Currently the only contract that is called when we open a pack is `CardFaucet`. The event generated from this contract is called `GeneratedCard`.

run this to generate ehtcontract
```
make oracle-abigen
```

Also, there is loomchainbackend which is a client for theh ethcontract. Some of the methods are not implemented yet because we don't use them. Please make sure you implement it if needed.

## Latest Plasma Block Number

We need to set latest plasma block number using `zb-cli` so that the oracle can start from the block where we already have the contract deployed. If not set, it will poll the data from block number 1

```
# start from block 200
./bin/zb-cli -k priv set_last_plasma_block_num -n 200
```


## Starting blocks

Dev: the valid block starts from block number 196492