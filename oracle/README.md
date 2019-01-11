# Gamechain Oracle

The code is mostly based on loomchain's [Oracle](https://github.com/loomnetwork/loomchain/tree/master/gateway). The upstream is Plasmachain which runs multiple Ethereum contracts using loomchain binary. The downstream is Gamechain which is a Go contracts.

The idea of this Oracle is to fetch all the events related to open pack and sycn to Gamechain. Ideally, we should have all of our cards we open on [loom.games](https://loom.games/) to Gamehchain to actually play the Zombie Battleground.

To fetch data from Plasmachain we need contract's abi to generate ethcontracts. Currently the only contract that is called when we open a pack is `CardFaucet`. The event generated from this contract is called `GeneratedCard`.

run this to generate ehtcontract
```
make oracle-abigen
```

Also, there is loomchainbackend which is a client for theh ethcontract. Some of the methods are not implemented yet because we don't use them. Please make sure you implement it if needed.