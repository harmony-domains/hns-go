// Copyright 2019 Weald Technology Trading
// Modified January 2023: John Whitton https://github.com/john_whitton
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
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jw-1ns/go-1ns/contracts/baseregistrar"
)

// BaseRegistrar is the structure for the registrar
type BaseRegistrar struct {
	backend      bind.ContractBackend
	domain       string
	Contract     *baseregistrar.Contract
	ContractAddr common.Address
}

// NewBaseRegistrar obtains the registrar contract for a given domain
func NewBaseRegistrar(backend bind.ContractBackend, domain string) (*BaseRegistrar, error) {
	address, err := RegistrarContractAddress(backend, domain)
	if err != nil {
		return nil, err
	}

	if address == UnknownAddress {
		return nil, fmt.Errorf("no registrar for domain %s", domain)
	}

	contract, err := baseregistrar.NewContract(address, backend)
	if err != nil {
		return nil, err
	}

	supported, err := contract.SupportsInterface(nil, [4]byte{0x28, 0xed, 0x4f, 0x6c})
	if err != nil {
		return nil, err
	}
	if !supported {
		return nil, fmt.Errorf("purported registrar for domain %s does not support reclaim functionality", domain)
	}

	return &BaseRegistrar{
		backend:      backend,
		domain:       domain,
		Contract:     contract,
		ContractAddr: address,
	}, nil
}

// RegisteredWith returns "permanent" or "none" for the
// registrar on which this name is registered
func (r *BaseRegistrar) RegisteredWith(domain string) (string, error) {
	// See if we're registered - fetch the owner to find out
	registry, err := NewRegistry(r.backend)
	if err != nil {
		return "", err
	}
	owner, err := registry.Owner(domain)
	if err != nil {
		return "", err
	}

	// No temporary registrar or no entry in same
	if owner == UnknownAddress {
		return "none", nil
	}
	return "permanent", nil
}

// Owner obtains the owner of the underlying token that represents the name.
func (r *BaseRegistrar) Owner(domain string) (common.Address, error) {
	name, err := UnqualifiedName(domain, r.domain)
	if err != nil {
		return UnknownAddress, err
	}
	labelHash, err := LabelHash(name)
	if err != nil {
		return UnknownAddress, err
	}
	owner, err := r.Contract.OwnerOf(nil, new(big.Int).SetBytes(labelHash[:]))
	// Registrar reverts rather than provide a 0 owner, so...
	if err != nil && err.Error() == "execution reverted" {
		return UnknownAddress, nil
	}
	return owner, err
}

// SetOwner sets the owner of the token holding the name
func (r *BaseRegistrar) SetOwner(opts *bind.TransactOpts, domain string, newOwner common.Address) (*types.Transaction, error) {
	name, err := UnqualifiedName(domain, r.domain)
	if err != nil {
		return nil, err
	}
	owner, err := r.Owner(name)
	if err != nil {
		return nil, err
	}
	labelHash, err := LabelHash(name)
	if err != nil {
		return nil, err
	}
	id := new(big.Int).SetBytes(labelHash[:])
	return r.Contract.TransferFrom(opts, owner, newOwner, id)
}

// Expiry obtains the unix timestamp at which the registration expires.
func (r *BaseRegistrar) Expiry(domain string) (*big.Int, error) {
	name, err := UnqualifiedName(domain, r.domain)
	if err != nil {
		return nil, err
	}
	labelHash, err := LabelHash(name)
	if err != nil {
		return nil, err
	}
	id := new(big.Int).SetBytes(labelHash[:])
	return r.Contract.NameExpires(nil, id)
}

// Reclaim reclaims a domain by the owner
func (r *BaseRegistrar) Reclaim(opts *bind.TransactOpts, domain string, newOwner common.Address) (*types.Transaction, error) {
	name, err := UnqualifiedName(domain, r.domain)
	if err != nil {
		return nil, err
	}
	labelHash, err := LabelHash(name)
	if err != nil {
		return nil, err
	}
	id := new(big.Int).SetBytes(labelHash[:])
	return r.Contract.Reclaim(opts, id, newOwner)
}
