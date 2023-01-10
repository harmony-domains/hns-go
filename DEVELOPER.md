# Overview

## Development Approach

For [coredns](https://github.com/coredns/coredns) we will build a pugin [coredns-hns](https://github.com/jw-ens-domains/coredns-hns) which interacts with [hns-go](https://github.com/jw-ens-domains/hns-go) leveraging code from [go-ens](https://github.com/wealdtech/go-ens) and [coredns-ens](https://github.com/wealdtech/coredns-ens). [hns-go](https://github.com/jw-ens-domains/hns-go) will have a focus of looking up dns entries using the dnsresolver and registry contracts. We will also create ethcontroller for testing registering of domains. Code not needed for retrieval of DNS entries will be omitted from the initial release, this includes IPFS support, auctions, deeds and other miscellaneous items.

## Geting Started 


Clone the repository

```
# Clone the repository
git clone https://github.com/jw-ens-domains/hns-go.git
cd hns-go

# update the dependencies
go get -u ./...
go mod tidy

# Build
go build

# Test
## Notes:
### Modify your eth_client as needed `client, _ := ethclient.Dial("https://ropsten.infura.io/v3/831a5442dc2e4536a9f8dee4ea1707a6")`
### If running locally will need to deploy your contracts using https://github.com/polymorpher/ens-deployer
### You may also need to create test data by registering domains

go test

# Run specific tests
go test -run name_test.go

```

## Development Tasks

### Updating contract abis

To update our contracts to the latest abi we can copy the abi from [ens-deployer](https://github.com/polymorpher/ens-deployer) and generate the contract.go file.

*Note: if adding a new contract we will also need to create a  contracts subdirectory with a `generate.go` file in the contracts subdirectory (e.g. dnsresolver folder)**

```
cd contracts
mkdir dnsresolver
//go:generate abigen -abi contract.abi -out contract.go -pkg dnsresolver -type Contract
```

Following is an example for updating `dnsresolver`

```
cd contracts/dnsresolver/
cp ../../../ens-deployer/contract/abi/PublicResolver.json contract.abi
go generate
```

After this we will need to add and update the corresponding calls to the contract. Typically there is a corresponding file such as `dnsresolver.go` which needs to be updated to reflect the changes we have just made to the contract (e.g. `contracts/dnsresolver.go`).

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





