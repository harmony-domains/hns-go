# Overview

## Development Approach

For [coredns](https://github.com/coredns/coredns) we will build a pugin [coredns-1ns](https://github.com/jw-1ns/coredns-1ns) which interacts with [go-1ns](https://github.com/jw-1ns/go-1ns) leveraging code from [go-ens](https://github.com/wealdtech/go-ens) and [coredns-ens](https://github.com/wealdtech/coredns-ens). [go-1ns](https://github.com/jw-1ns/go-1ns) will have a focus of looking up dns entries using the PubliceRsolver and Registry contracts. We will also create Registrarcontroller for testing registering of domains. Code not needed for retrieval of DNS entries will be omitted from the initial release, this includes IPFS support, auctions, deeds and other miscellaneous items.

## Geting Started

Clone the repository

```bash
# Clone the repository
git clone https://github.com/jw-1ns/go-1ns.git
cd go-1ns

# update the dependencies
go get -u ./...
go mod tidy

# Build
go build

```

Testing

You will need to run a local ganche instance and deploy ens contracts and register sample domains for the tests to work

This will [register the following domains](https://github.com/jw-1ns/ens-deployer/blob/main/contract/deploy/dnsSample.ts#LL133-L138C54)

* await registerDomain('test', alice, '128.0.0.1')
* await registerDomain('testa', alice, '128.0.0.2')
* await registerDomain('testb', bob, '128.0.0.3')
* await registerDomain('testlongdomain', bob, '128.0.0.4')
* await registerDomain('testxyz', alice, '128.0.0.5')

```bash
git clone https://github.com/jw-1ns/ens-deployer.git
cd ens-deployer
# git checkout 1ns

# start local ganache in window 1
cd env
./ganache-new.sh

# deploy the contracts in window 2
cd contracts
yarn install
cp .env.sample .env
yarn deploy-dns
```

Running go tests

We explicity name the tests we are running using the Name of the test. *Note: this means when writing test each functional area is held within one test file and each test should have a unique prefix.*

To run all tests matching a prefix `go test -run TestName`
To match a specific test `go test -run "^TestName$"`

```bash

# Run all tests
go test

# Test configuration i.e. local ganache contracts are deployed correctly and expected domains have been registered

go test -run TestConfig

# We do not use IPFS so no need to test contenthash
# go test -run TestContenthash

# DNS registration is not currently supported  by go-1ns (we use ens-deployer currently) so no need to test dnsresolver
# go test -run TestA

# Test miscellaneous tests
go test -run "^TestNormaliseDomain"
go test -run TestNormaliseDomain
go test -run TestNormaliseDomainStrict
go test -run TestTld
go test -run TestDomainPart
go test -run TestUnqualifiedName

# Test name 
# go test -run "^TestName$"
go test -run TestName

# Test namehash
go test -run TestNameHash
go test -run TestLabelHash

# Test registry
go test -run TestResolve

```

## Development Tasks

### Updating contract abis

To update our contracts to the latest abi we can copy the abi from [ens-deployer](https://github.com/polymorpher/ens-deployer) and generate the contract.go file.

Prerequiste:  [abigen](https://geth.ethereum.org/docs/tools/abigen) is needed, this can be installed by installing [geth](https://geth.ethereum.org/docs/getting-started/installing-geth) for example

```bash
brew tap ethereum/ethereum
brew install ethereum
```

*Note: if adding a new contract we will also need to create a  contracts subdirectory with a `generate.go` file in the contracts subdirectory (e.g. dnsresolver folder)**

```bash
cd contracts
mkdir publicresolver
//go:generate abigen -abi contract.abi -out contract.go -pkg publicresolver -type Contract
```

Following is an example for updating `dnsresolver`

```bash
cd contracts/publicresolver/
cp ../../../ens-deployer/contract/abi/PublicResolver.json contract.abi
go generate
```

After this we will need to add and update the corresponding calls to the contract. Typically there is a corresponding file such as `publicresolver.go` which needs to be updated to reflect the changes we have just made to the contract (e.g. `contracts/publicresolver.go`).

We then need to build and fix any compilation errors.

```bash
cd ../..
go build
```

Then we can run your tests again

```bash
go test -run name_test.go
```

This generates `contract.go`

Additional references

* [Go Contract Bindings](https://geth.ethereum.org/docs/developers/dapp-developer/native-bindings) ; Geth ethereum org
* [Smart Contract Compilation $ ABI](https://goethereumbook.org/smart-contract-compile/): Go Ethereum Book

## View go docs

Optionally install godoc

`go install golang.org/x/tools/cmd/godoc@latest`

To run the godoc server

`godoc -http=:6061 -index`

and view the docs at

[http://localhost:6061/pkg/](http://localhost:6061/pkg/)

view go-1ns docs at

[http://localhost:6061/pkg/github.com/jw-1ns/go-1ns/](http://localhost:6061/pkg/github.com/jw-1ns/go-1ns/)

or specifically
[http://localhost:6061/pkg/github.com/jw-1ns/go-1ns/](http://localhost:6061/pkg/github.com/jw-1ns/go-1ns/)

## Publishing a module

See [here](https://go.dev/doc/modules/publishing) and [here](https://go.dev/blog/publishing-go-modules).

```bash
git tag v0.1.5
git push origin v0.1.5
```
