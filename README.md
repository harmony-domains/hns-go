# go-1ns

[![Tag](https://img.shields.io/github/tag/jw-1ns/go-1ns.svg)](https://github.com/jw-1ns/go-1ns/releases/)
[![License](https://img.shields.io/github/license/jw-1ns/go-1ns.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/jw-1ns/go-1ns?status.svg)](https://godoc.org/github.com/jw-1ns/go-1ns)
[![Travis CI](https://img.shields.io/travis/jw-1ns/go-1ns.svg)](https://travis-ci.org/jw-1ns/go-1ns)
[![codecov.io](https://img.shields.io/codecov/c/github/jw-1ns/go-1ns.svg)](https://codecov.io/github/jw-1ns/go-1ns)
[![Go Report Card](https://goreportcard.com/badge/github.com/jw-1ns/go-1ns)](https://goreportcard.com/report/github.com/jw-1ns/go-1ns)

Go module to simplify interacting with the [1 Name Service](https://1ns.domains/) contracts.
Initial version copied from [go-ens](https://github.com/wealdtech/go-ens)


## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Install

`go-1ns` is a standard Go module which can be installed with:

```sh
go get github.com/jw-1ns/go-1ns
```

## Usage

`go-1ns` provides simple access to the [1 Name Service](https://1ns.domains/) (1ns) contracts.

### Resolution

The most commonly-used feature of 1ns is resolution: converting an 1ns name to an Ethereum address.  `go-1ns` provides a simple call to allow this:

```go
address, err := onens.Resolve(client, domain)
```

where `client` is a connection to an Harmony client and `domain` is the fully-qualified name you wish to resolve (e.g. `foo.mydomain.country`) (full examples for using this are given in the [Example](#Example) section below).

The reverse process, converting an address to an 1ns name, is just as simple:

```go
domain, err := onens.ReverseResolve(client, address)
```

Note that if the address does not have a reverse resolution this will return "".  If you just want a string version of an address for on-screen display then you can use `onens.Format()`, for example:

```go
fmt.Printf("The address is %s\n", onens.Format(client, address))
```

This will carry out reverse resolution of the address and print the name if present; if not it will print a formatted version of the address.


### Management of names

A top-level name is one that sits directly underneath `.country`, for example `mydomain.country`.  Lower-level names, such as `foo.mydomain.country` are covered in the following section.  `go-1ns` provides a simplified `Name` interface to manage top-level, removing the requirement to understand registrars, controllers, _etc._

Starting out with names in `go-1ns` is easy:

```go
client, err := ethclient.Dial("https://api.s0.t.hmny.io")
name, err := onens.NewName(client, "mydomain.country")
```

Addresses can be set and obtained using the address functions, for example to get an address:

```go
COIN_TYPE_ETHEREUM := uint64(60)
address, err := name.Address(COIN_TYPE_ETHEREUM)
```

1ns supports addresses for multiple coin types; values of coin types can be found at https://github.com/satoshilabs/slips/blob/master/slip-0044.md

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
	onens "github.com/jw-1ns/go-1ns"
)

func main() {
	// Replace SECRET with your own access token for this example to work.
	client, err := ethclient.Dial("https://api.s0.t.hmny.io")
	if err != nil {
		panic(err)
	}

	// Resolve a name to an address.
	domain := "test.country"
	address, err := onens.Resolve(client, domain)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Address of %s is %s\n", domain, address.Hex())

	// Reverse resolve an address to a name.
	reverse, err := onens.ReverseResolve(client, address)
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

Contributions welcome. Please check out [the issues](https://github.com/jw-1ns/go-1ns/issues).

## License

[Apache-2.0](LICENSE) © 2022 John Whitton
