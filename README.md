# hns-go

[![Tag](https://img.shields.io/github/tag/harmony-domains/hns-go.svg)](https://github.com/harmony-domains/hns-go/releases/)
[![License](https://img.shields.io/github/license/harmony-domains/hns-go.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/harmony-domains/hns-go?status.svg)](https://godoc.org/github.com/harmony-domains/hns-go)
[![Travis CI](https://img.shields.io/travis/harmony-domains/hns-go.svg)](https://travis-ci.org/harmony-domains/hns-go)
[![codecov.io](https://img.shields.io/codecov/c/github/harmony-domains/hns-go.svg)](https://codecov.io/github/harmony-domains/hns-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/harmony-domains/hns-go)](https://goreportcard.com/report/github.com/harmony-domains/hns-go)

Go module to simplify interacting with the [Harmony Name Service](https://hns.domains/) contracts.
Initial version copied from [go-ens](https://github.com/wealdtech/go-ens)


## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Install

`hns-go` is a standard Go module which can be installed with:

```sh
go get github.com/harmony-domains/hns-go
```

## Usage

`hns-go` provides simple access to the [Harmony Name Service](https://hns.domains/) (HNS) contracts.

### Resolution

The most commonly-used feature of HNS is resolution: converting an HNS name to an Ethereum address.  `hns-go` provides a simple call to allow this:

```go
address, err := hns.Resolve(client, domain)
```

where `client` is a connection to an Ethereum client and `domain` is the fully-qualified name you wish to resolve (e.g. `foo.mydomain.eth`) (full examples for using this are given in the [Example](#Example) section below).

The reverse process, converting an address to an HNS name, is just as simple:

```go
domain, err := hns.ReverseResolve(client, address)
```

Note that if the address does not have a reverse resolution this will return "".  If you just want a string version of an address for on-screen display then you can use `hns.Format()`, for example:

```go
fmt.Printf("The address is %s\n", hns.Format(client, address))
```

This will carry out reverse resolution of the address and print the name if present; if not it will print a formatted version of the address.


### Management of names

A top-level name is one that sits directly underneath `.eth`, for example `mydomain.eth`.  Lower-level names, such as `foo.mydomain.eth` are covered in the following section.  `hns-go` provides a simplified `Name` interface to manage top-level, removing the requirement to understand registrars, controllers, _etc._

Starting out with names in `hns-go` is easy:

```go
client, err := ethclient.Dial("https://infura.io/v3/SECRET")
name, err := hns.NewName(client, "mydomain.eth")
```

Addresses can be set and obtained using the address functions, for example to get an address:

```go
COIN_TYPE_ETHEREUM := uint64(60)
address, err := name.Address(COIN_TYPE_ETHEREUM)
```

HNS supports addresses for multiple coin types; values of coin types can be found at https://github.com/satoshilabs/slips/blob/master/slip-0044.md

### Registering and extending names

Most operations on a domain will involve setting resolvers and resolver information.


### Management of subdomains

Because subdomains have their own registrars they do not work with the `Name` interface.

### Example

```go
package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	hns "github.com/harmony-domains/hns-go"
)

func main() {
	// Replace SECRET with your own access token for this example to work.
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/SECRET")
	if err != nil {
		panic(err)
	}

	// Resolve a name to an address.
	domain := "ethereum.eth"
	address, err := hns.Resolve(client, domain)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Address of %s is %s\n", domain, address.Hex())

	// Reverse resolve an address to a name.
	reverse, err := hns.ReverseResolve(client, address)
	if err != nil {
		panic(err)
	}
	if reverse == "" {
		fmt.Printf("%s has no reverse lookup\n", address.Hex())
	} else {
		fmt.Printf("Name of %s is %s\n", address.Hex(), reverse)
	}
}
```

## Maintainers

John Whitton: [@john_whitton](https://github.com/john_whitton).

## Contribute

Contributions welcome. Please check out [the issues](https://github.com/harmony-domains/hns-go/issues).

## License

[Apache-2.0](LICENSE) © 2022 Harmony Domains
