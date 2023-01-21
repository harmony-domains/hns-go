# Overview

## Development Approach

For [coredns](https://github.com/coredns/coredns) we will build a pugin [coredns-1ns](https://github.com/jw-ens-domains/coredns-1ns) which interacts with [go-1ns](https://github.com/jw-ens-domains/go-1ns) leveraging code from [go-ens](https://github.com/wealdtech/go-ens) and [coredns-ens](https://github.com/wealdtech/coredns-ens). [go-1ns](https://github.com/jw-ens-domains/go-1ns) will have a focus of looking up dns entries using the PubliceRsolver and Registry contracts. We will also create Registrarcontroller for testing registering of domains. Code not needed for retrieval of DNS entries will be omitted from the initial release, this includes IPFS support, auctions, deeds and other miscellaneous items.

## Geting Started 


Clone the repository

```
# Clone the repository
git clone https://github.com/jw-ens-domains/go-1ns.git
cd go-1ns

# update the dependencies
go get -u ./...
go mod tidy

# Build
go build

# Test
## Notes:
### Modify your eth_client as needed `client, _ := ethclient.Dial("http://localhost:8545")`
### If running locally will need to deploy your contracts using https://github.com/polymorpher/ens-deployer
### You may also need to create test data by registering domains

go test

```

## Development Tasks

### Updating contract abis

To update our contracts to the latest abi we can copy the abi from [ens-deployer](https://github.com/polymorpher/ens-deployer) and generate the contract.go file.

*Note: if adding a new contract we will also need to create a  contracts subdirectory with a `generate.go` file in the contracts subdirectory (e.g. dnsresolver folder)**

```
cd contracts
mkdir publicresolver
//go:generate abigen -abi contract.abi -out contract.go -pkg publicresolver -type Contract
```

Following is an example for updating `dnsresolver`

```
cd contracts/publicresolver/
cp ../../../ens-deployer/contract/abi/PublicResolver.json contract.abi
go generate
```

After this we will need to add and update the corresponding calls to the contract. Typically there is a corresponding file such as `publicresolver.go` which needs to be updated to reflect the changes we have just made to the contract (e.g. `contracts/publicresolver.go`).

We then need to build and fix any compilation errors.

```
cd ../..
go build
```

Then we can run your tests again

```
go test -run name_test.go
```

This generates `contract.go` 

Additional references
* [Go Contract Bindings](https://geth.ethereum.org/docs/developers/dapp-developer/native-bindings) ; Geth ethereum org
* [Smart Contract Compilation $ ABI](https://goethereumbook.org/smart-contract-compile/): Go Ethereum Book 

## View go docs

To run the godoc server

`godoc -http=:6061 -index`

and view the docs at

[http://localhost:6061/pkg/](http://localhost:6061/pkg/)

view go-1ns docs at 

[http://localhost:6061/pkg/github.com/jw-1ns/go-1ns/](http://localhost:6061/pkg/github.com/jw-1ns/go-1ns/)


## Publishing a module

See [here](https://go.dev/doc/modules/publishing) and [here](https://go.dev/blog/publishing-go-modules).

```
git tag v1.0.0
git push origin v1.0.0
```




