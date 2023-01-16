// Copyright 2017 Weald Technology Trading
// Modified December 2022: John Whitton https://github.com/john_whitton
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

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/harmony-domains/hns-go/contracts/auctionregistrar"
	"github.com/harmony-domains/hns-go/contracts/registry"
	"github.com/jw-1ns/go-1ns/contracts/ensregistry"
	"github.com/jw-1ns/go-1ns/util"
	"github.com/pkg/errors"
)

// Registry is the structure for the ensregistry contract
type Registry struct {
	backend      bind.ContractBackend
	Contract     *ensregistry.Contract
	ContractAddr common.Address
}

// NewRegistry obtains the ENS registry
func NewRegistry(backend bind.ContractBackend) (*Registry, error) {
	contract, err := ensregistry.NewContract(config.ENSRegistry, backend)
	if err != nil {
		return nil, err
	}
	return &Registry{
		backend:      backend,
		Contract:     contract,
		ContractAddr: config.ENSRegistry,
	}, nil
}

// Owner returns the address of the owner of a name
func (r *Registry) Owner(name string) (common.Address, error) {
	nameHash, err := NameHash(name)
	if err != nil {
		return UnknownAddress, err
	}
	return r.Contract.Owner(nil, nameHash)
}

// ResolverAddress returns the address of the resolver for a name
func (r *Registry) ResolverAddress(name string) (common.Address, error) {
	nameHash, err := NameHash(name)
	fmt.Printf("ResolverAddress Name: %+v\n", name)
	fmt.Printf("ResolverAddress Name Hash: %+v\n", nameHash)
	if err != nil {
		return UnknownAddress, err
	}
	return r.Contract.Resolver(nil, nameHash)
}

// SetResolver sets the resolver for a name
func (r *Registry) SetResolver(opts *bind.TransactOpts, name string, address common.Address) (*types.Transaction, error) {
	nameHash, err := NameHash(name)
	if err != nil {
		return nil, err
	}
	return r.Contract.SetResolver(opts, nameHash, address)
}

// Resolver returns the resolver for a name
func (r *Registry) Resolver(name string) (*Resolver, error) {
	address, err := r.ResolverAddress(name)
	if err != nil {
		return nil, err
	}
	fmt.Printf("RegistryResolverAddress%+v\n", address)
	return NewResolverAt(r.backend, name, address)
}

// SetOwner sets the ownership of a domain
func (r *Registry) SetOwner(opts *bind.TransactOpts, name string, address common.Address) (*types.Transaction, error) {
	nameHash, err := NameHash(name)
	if err != nil {
		return nil, err
	}
	return r.Contract.SetOwner(opts, nameHash, address)
}

// SetSubdomainOwner sets the ownership of a subdomain, potentially creating it in the process
func (r *Registry) SetSubdomainOwner(opts *bind.TransactOpts, name string, subname string, address common.Address) (*types.Transaction, error) {
	nameHash, err := NameHash(name)
	if err != nil {
		return nil, err
	}
	labelHash, err := LabelHash(subname)
	if err != nil {
		return nil, err
	}
	return r.Contract.SetSubnodeOwner(opts, nameHash, labelHash, address)
}

// RegistryContractAddress obtains the address of the registry contract for a chain.
// Get the Registry contract address from config
func RegistryContractAddress(backend bind.ContractBackend) (common.Address, error) {
	// config := getConfig()
	return config.ENSRegistry, nil
}

// RegistryContractFromRegistrar obtains the registry contract given an
// existing registrar contract
func RegistryContractFromRegistrar(backend bind.ContractBackend, registrar *auctionregistrar.Contract) (*registry.Contract, error) {
	if registrar == nil {
		return nil, errors.New("no registrar contract")
	}
	registryAddress, err := registrar.Ens(nil)
	if err != nil {
		return nil, err
	}
	return registry.NewContract(registryAddress, backend)
}

// SetResolver sets the resolver for a name
func SetResolver(session *ensregistry.ContractSession, name string, resolverAddr *common.Address) (*types.Transaction, error) {
	nameHash, err := NameHash(name)
	if err != nil {
		return nil, err
	}
	return session.SetResolver(nameHash, *resolverAddr)
}

// SetSubdomainOwner sets the owner for a subdomain of a name
func SetSubdomainOwner(session *ensregistry.ContractSession, name string, subdomain string, ownerAddr *common.Address) (*types.Transaction, error) {
	nameHash, err := NameHash(name)
	if err != nil {
		return nil, err
	}
	labelHash, err := LabelHash(subdomain)
	if err != nil {
		return nil, err
	}
	return session.SetSubnodeOwner(nameHash, labelHash, *ownerAddr)
}

// CreateRegistrySession creates a session suitable for multiple calls
func CreateRegistrySession(chainID *big.Int, wallet *accounts.Wallet, account *accounts.Account, passphrase string, contract *ensregistry.Contract, gasPrice *big.Int) *ensregistry.ContractSession {
	// Create a signer
	signer := util.AccountSigner(chainID, wallet, account, passphrase)

	// Return our session
	session := &ensregistry.ContractSession{
		Contract: contract,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:     account.Address,
			Signer:   signer,
			GasPrice: gasPrice,
		},
	}

	return session
}
