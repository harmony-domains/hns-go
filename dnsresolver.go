// Copyright 2017-2019 Weald Technology Trading
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package onens

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jw-1ns/go-1ns/contracts/publicresolver"
	"golang.org/x/crypto/sha3"
)

// DNSResolver is the structure for the DNS resolver contract
type DNSResolver struct {
	backend      bind.ContractBackend
	domain       string
	Contract     *publicresolver.Contract
	ContractAddr common.Address
}

// NewDNSResolver creates a new DNS resolver for a given domain
func NewDNSResolver(backend bind.ContractBackend, domain string) (*DNSResolver, error) {
	registry, err := NewRegistry(backend)
	if err != nil {
		return nil, err
	}
	address, err := registry.ResolverAddress(domain)
	fmt.Printf("NEWDNSResolver address: %+v\n", address)
	if err != nil {
		return nil, err
	}

	return NewDNSResolverAt(backend, domain, address)
}

// NewDNSResolverAt creates a new DNS resolver for a given domain at a given address
func NewDNSResolverAt(backend bind.ContractBackend, domain string, address common.Address) (*DNSResolver, error) {
	contract, err := publicresolver.NewContract(address, backend)
	if err != nil {
		return nil, err
	}

	// Ensure that this is a DNS resolver
	supported, err := contract.SupportsInterface(nil, [4]byte{0xa8, 0xfa, 0x56, 0x82})
	if err != nil {
		return nil, err
	}
	if !supported {
		return nil, fmt.Errorf("%s is not a DNS resolver contract", address.Hex())
	}
	fmt.Println("NEWDNSResolverAt is supported")

	return &DNSResolver{
		backend:      backend,
		domain:       domain,
		Contract:     contract,
		ContractAddr: address,
	}, nil
}

// Record obtains an RRSet for a name
func (r *DNSResolver) Record(name string, rrType uint16) ([]byte, error) {
	nameHash, err := NameHash(r.domain)
	fmt.Printf("name: %+v\n", name)
	fmt.Printf("r.domain: %+v\n", r.domain)
	fmt.Printf("nameHash: %+v\n", nameHash)
	fmt.Printf("DNSWireFormatDomainHash(name): %+v\n", DNSWireFormatDomainHash(name))
	if err != nil {
		return nil, err
	}
	return r.Contract.DnsRecord(nil, nameHash, DNSWireFormatDomainHash(name), rrType)
}

// HasRecords returns true if the given name has any RRsets
func (r *DNSResolver) HasRecords(name string) (bool, error) {
	nameHash, err := NameHash(r.domain)
	if err != nil {
		return false, err
	}
	return r.Contract.HasDNSRecords(nil, nameHash, DNSWireFormatDomainHash(name))
}

// SetRecords sets one or more RRSets
func (r *DNSResolver) SetRecords(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	nameHash, err := NameHash(r.domain)
	// nameHashBack := []byte(string([]byte(nameHash[:])))
	// nameHashString := string([]byte(nameHash[:]))
	nameHashHex := hex.EncodeToString([]byte(nameHash[:]))
	nameHashBack2, err := hex.DecodeString(nameHashHex)
	// nameHashHexStringBack := string(nameHashBack2)
	fmt.Printf("r.domain: %+v\n", r.domain)
	fmt.Printf("nameHash: %+v\n", nameHash)
	// fmt.Printf("nameHash back: %v\n", nameHashBack)
	// fmt.Printf("nameHash string: %+v\n", nameHashString)
	fmt.Printf("nameHash hex: %+v\n", nameHashHex)
	// fmt.Printf("nameHash hex back: %+v\n", nameHashHexStringBack)
	fmt.Printf("nameHash back2: %v\n", nameHashBack2)

	// existingRec, err := r.Contract.DnsRecord(nil, nameHash, DNSWireFormatDomainHash("a.test.country."), 1)
	// if err != nil {
	// 	fmt.Printf("err: %+v\n", err)
	// }
	// fmt.Printf("existingRec: %+v\n", existingRec)
	if err != nil {
		return nil, err
	}
	// return r.Contract.SetDNSRecords(opts, nameHash2, data)
	return r.Contract.SetDNSRecords(opts, nameHash, data)
}

// ClearRecords clears all records for a domain
func (r *DNSResolver) ClearRecords(opts *bind.TransactOpts) (*types.Transaction, error) {
	nameHash, err := NameHash(r.domain)
	// nameHash, err := NameHash("badtest.country")
	if err != nil {
		return nil, err
	}
	return r.Contract.ClearRecords(opts, nameHash)
}

// Zonehash returns the zone hash of the domain
func (r *DNSResolver) Zonehash() ([]byte, error) {
	nameHash, err := NameHash(r.domain)
	if err != nil {
		return nil, err
	}
	return r.Contract.Zonehash(nil, nameHash)
}

// SetZonehash sets the zone hash of the domain
func (r *DNSResolver) SetZonehash(opts *bind.TransactOpts, zonehash []byte) (*types.Transaction, error) {
	nameHash, err := NameHash(r.domain)
	if err != nil {
		return nil, err
	}
	return r.Contract.SetZonehash(opts, nameHash, zonehash)
}

// DNSWireFormatDomainHash hashes a domain name in wire format
func DNSWireFormatDomainHash(domain string) (hash [32]byte) {
	sha := sha3.NewLegacyKeccak256()
	// //nolint:golint,errcheck
	sha.Write(DNSWireFormat(domain))
	sha.Sum(hash[:0])
	return
}

// DNSWireFormat turns a domain name in to wire format
func DNSWireFormat(domain string) []byte {
	// Remove leading and trailing dots
	domain = strings.TrimLeft(domain, ".")
	domain = strings.TrimRight(domain, ".")
	domain = strings.ToLower(domain)

	if domain == "" {
		return []byte{0x00}
	}

	bytes := make([]byte, len(domain)+2)
	pieces := strings.Split(domain, ".")
	offset := 0
	for _, piece := range pieces {
		bytes[offset] = byte(len(piece))
		offset++
		copy(bytes[offset:offset+len(piece)], piece)
		offset += len(piece)
	}
	return bytes
}
