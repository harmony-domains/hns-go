// Copyright 2017-2019 Weald Technology Trading
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

/*
Following is a list of errors, events and public functions for registrarcontroller
[
  "constructor(address,address,uint256,uint256,address,address,bytes32,string)",
  "error CommitmentTooNew(bytes32)",
  "error CommitmentTooOld(bytes32)",
  "error DurationTooShort(uint256)",
  "error InsufficientValue()",
  "error MaxCommitmentAgeTooHigh()",
  "error MaxCommitmentAgeTooLow()",
  "error NameNotAvailable(string)",
  "error ResolverRequiredWhenDataSupplied()",
  "error Unauthorised(bytes32)",
  "error UnexpiredCommitmentExists(bytes32)",
  "event NameRegistered(string,bytes32 indexed,address indexed,uint256,uint256,uint256)",
  "event NameRenewed(string,bytes32 indexed,uint256,uint256)",
  "event OwnershipTransferred(address indexed,address indexed)",
  "function MIN_REGISTRATION_DURATION() view returns (uint256)",
  "function available(string) view returns (bool)",
  "function baseExtension() view returns (string)",
  "function baseNode() view returns (bytes32)",
  "function commit(bytes32)",
  "function commitments(bytes32) view returns (uint256)",
  "function makeCommitment(string,address,uint256,bytes32,address,bytes[],bool,uint32,uint64) pure returns (bytes32)",
  "function maxCommitmentAge() view returns (uint256)",
  "function minCommitmentAge() view returns (uint256)",
  "function nameWrapper() view returns (address)",
  "function owner() view returns (address)",
  "function prices() view returns (address)",
  "function recoverFunds(address,address,uint256)",
  "function register(string,address,uint256,bytes32,address,bytes[],bool,uint32,uint64) payable",
  "function renew(string,uint256) payable",
  "function renewWithFuses(string,uint256,uint32,uint64) payable",
  "function renounceOwnership()",
  "function rentPrice(string,uint256) view returns (tuple(uint256,uint256))",
  "function reverseRegistrar() view returns (address)",
  "function supportsInterface(bytes4) pure returns (bool)",
  "function transferOwnership(address)",
  "function valid(string) pure returns (bool)",
  "function withdraw()"
]

*/

package onens

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jw-1ns/go-1ns/contracts/registrarcontroller"
	"github.com/pkg/errors"
)

// RegistrarController is the structure for the registrar controller contract
type RegistrarController struct {
	backend      bind.ContractBackend
	Contract     *registrarcontroller.Contract
	ContractAddr common.Address
	domain       string
}

// NewRegistrarController creates a new controller for a given domain
func NewRegistrarController(backend bind.ContractBackend, domain string) (*RegistrarController, error) {
	registry, err := NewRegistry(backend)
	if err != nil {
		return nil, err
	}
	resolver, err := registry.Resolver(domain)
	if err != nil {
		return nil, err
	}

	// Obtain the controller from the resolver
	// permanentRegistrar: '0x018fac06'
	controllerAddress, err := resolver.InterfaceImplementer([4]byte{0x01, 0x8f, 0xac, 0x06})
	if err != nil {
		return nil, err
	}

	return NewRegistrarControllerAt(backend, domain, controllerAddress)
}

// NewETHControllerAt creates a .eth controller at a given address
func NewRegistrarControllerAt(backend bind.ContractBackend, domain string, address common.Address) (*RegistrarController, error) {
	contract, err := registrarcontroller.NewContract(address, backend)
	if err != nil {
		return nil, err
	}
	return &RegistrarController{
		backend:      backend,
		Contract:     contract,
		ContractAddr: address,
		domain:       domain,
	}, nil
}

//	"function MIN_REGISTRATION_DURATION() view returns (uint256)",
//
// MinRegistrationDuration returns the minimum duration for which a name can be registered
func (c *RegistrarController) MinRegistrationDuration() (time.Duration, error) {
	tmp, err := c.Contract.MINREGISTRATIONDURATION(nil)
	if err != nil {
		return 0 * time.Second, err
	}

	return time.Duration(tmp.Int64()) * time.Second, nil
}

// "function available(string) view returns (bool)",
// IsAvailable returns true if the domain is available for registration.
func (c *RegistrarController) IsAvailable(domain string) (bool, error) {
	name, err := UnqualifiedName(domain, c.domain)
	if err != nil {
		return false, fmt.Errorf("invalid name %s", domain)
	}
	return c.Contract.Available(nil, name)
}

// "function baseExtension() view returns (string)",
// BaseExtension retrieves the baseExtension
func (c *RegistrarController) Basextension() (string, error) {
	return c.Contract.BaseExtension(nil)
}

// "function baseNode() view returns (bytes32)",
// BaseNode retreives the base node
func (c *RegistrarController) BaseNode(opts *bind.CallOpts) ([32]byte, error) {
	return c.Contract.BaseNode(opts)
}

// "function commit(bytes32)",
// // Commit sends a commitment to register a domain.
// func (c *RegistrarController) Commit(opts *bind.TransactOpts, commitment [32]byte) (*types.Transaction, error) {
// 	return c.Contract.Commit(opts, commitment)
// }

// Commit sends a commitment to register a domain.
// func (_Contract *ContractCaller) MakeCommitment(opts *bind.CallOpts, name string, owner common.Address, duration *big.Int, secret [32]byte, resolver common.Address, data [][]byte, reverseRecord bool, fuses uint32, wrapperExpiry uint64) ([32]byte, error) {
func (c *RegistrarController) Commit(opts *bind.TransactOpts, domain string, owner common.Address, duration *big.Int, secret [32]byte) (*types.Transaction, error) {
	name, err := UnqualifiedName(domain, c.domain)
	if err != nil {
		return nil, fmt.Errorf("invalid name %s", domain)
	}

	// func (_Contract *ContractCaller) MakeCommitment(opts *bind.CallOpts, name string, owner common.Address, duration *big.Int, secret [32]byte, resolver common.Address, data [][]byte, reverseRecord bool, fuses uint32, wrapperExpiry uint64) ([32]byte, error) {
	// commitment, err := c.Contract.MakeCommitment(nil, name, owner, secret)
	config := getConfig()
	commitment, err := c.Contract.MakeCommitment(nil, name, owner, duration, secret, config.commitmentData.publicResover, config.commitmentData.calldata, config.commitmentData.reverseRecord, config.commitmentData.fuses, config.commitmentData.wrapperExpiry)
	if err != nil {
		return nil, errors.New("failed to create commitment")
	}

	if opts.Value != nil && opts.Value.Cmp(big.NewInt(0)) != 0 {
		return nil, errors.New("commitment should have 0 value")
	}

	return c.Contract.Commit(opts, commitment)
}

// "function commitments(bytes32) view returns (uint256)"
// Commitments returns the block timestamp of the commmitment
func (c *RegistrarController) Commitments(opts *bind.CallOpts, commitment [32]byte) (*big.Int, error) {
	return c.Contract.Commitments(opts, commitment)
}

// CommitmentTime states the time at which a commitment was registered on the blockchain.
func (c *RegistrarController) CommitmentTime(domain string, owner common.Address, duration *big.Int, secret [32]byte) (*big.Int, error) {
	hash, err := c.CommitmentHash(domain, owner, duration, secret)
	if err != nil {
		return nil, err
	}

	return c.Contract.Commitments(nil, hash)
}

// "function makeCommitment(string,address,uint256,bytes32,address,bytes[],bool,uint32,uint64) pure returns (bytes32)",
// CommitmentHash returns the commitment hash for a label/owner/secret tuple
func (c *RegistrarController) CommitmentHash(domain string, owner common.Address, duration *big.Int, secret [32]byte) (common.Hash, error) {
	name, err := UnqualifiedName(domain, c.domain)
	if err != nil {
		return common.BytesToHash([]byte{}), fmt.Errorf("invalid name %s", domain)
	}

	config := getConfig()
	commitment, err := c.Contract.MakeCommitment(nil, name, owner, duration, secret, config.commitmentData.publicResover, config.commitmentData.calldata, config.commitmentData.reverseRecord, config.commitmentData.fuses, config.commitmentData.wrapperExpiry)
	// commitment, err := c.Contract.MakeCommitment(nil, name, owner, duration, secret, resolver, data, reverseRecord, fuses, wrapperExpiry)
	if err != nil {
		return common.BytesToHash([]byte{}), err
	}
	hash := common.BytesToHash(commitment[:])
	return hash, err
}

// "function maxCommitmentAge() view returns (uint256)",
func (c *RegistrarController) MaxCommitmentInterval() (*big.Int, error) {
	return c.Contract.MaxCommitmentAge(nil)
}

// "function minCommitmentAge() view returns (uint256)",
// MinCommitmentInterval returns the minimum time that has to pass between a commit and reveal
func (c *RegistrarController) MinCommitmentInterval() (*big.Int, error) {
	return c.Contract.MinCommitmentAge(nil)
}

// "function nameWrapper() view returns (address)",
func (c *RegistrarController) NameWrapper() (common.Address, error) {
	return c.Contract.NameWrapper(nil)
}

// "function owner() view returns (address)",
func (c *RegistrarController) Owner() (common.Address, error) {
	return c.Contract.Owner(nil)
}

// "function prices() view returns (address)",
func (c *RegistrarController) Prices() (common.Address, error) {
	return c.Contract.Prices(nil)
}

// "function recoverFunds(address,address,uint256)",
func (c *RegistrarController) RecoverFunds(opts *bind.TransactOpts, token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return c.Contract.RecoverFunds(opts, token, to, amount)
}

// "function register(string,address,uint256,bytes32,address,bytes[],bool,uint32,uint64) payable",
func (c *RegistrarController) Register(opts *bind.TransactOpts, name string, owner common.Address, duration *big.Int, secret [32]byte, resolver common.Address, data [][]byte, reverseRecord bool, fuses uint32, wrapperExpiry uint64) (*types.Transaction, error) {
	return c.Contract.Register(opts, name, owner, duration, secret, resolver, data, reverseRecord, fuses, wrapperExpiry)
}

// "function renew(string,uint256) payable",
// func (c *RegistrarController) Renew(opts *bind.TransactOpts, name string, duration *big.Int) (*types.Transaction, error) {
// 	return c.Contract.Renew(opts, name, duration)
// }

// "function renewWithFuses(string,uint256,uint32,uint64) payable",
func (c *RegistrarController) RenewWithFuses(opts *bind.TransactOpts, name string, duration *big.Int, fuses uint32, wrapperExpiry uint64) (*types.Transaction, error) {
	return c.Contract.RenewWithFuses(opts, name, duration, fuses, wrapperExpiry)
}

// "function renounceOwnership()",
func (c *RegistrarController) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return c.Contract.RenounceOwnership(opts)
}

// "function rentPrice(string,uint256) view returns (tuple(uint256,uint256))",
func (c *RegistrarController) RentPrice(opts *bind.CallOpts, name string, duration *big.Int) (registrarcontroller.IPriceOraclePrice, error) {
	return c.Contract.RentPrice(opts, name, duration)
}

// RentCost returns the cost of rent in wei-per-second.
func (c *RegistrarController) RentCost(domain string) (*big.Int, error) {
	name, err := UnqualifiedName(domain, c.domain)
	if err != nil {
		return nil, fmt.Errorf("invalid name %s", domain)
	}
	//TODO Modify this to read the value returned from the PriceOracle type IPriceOraclePrice
	priceOraclePrice, err := c.Contract.RentPrice(nil, name, big.NewInt(1))
	if err != nil {
		return nil, fmt.Errorf("invalid price for %s", domain)
	}
	var price *big.Int = big.NewInt(0)
	price.Add(priceOraclePrice.Base, priceOraclePrice.Premium)
	return price, nil
}

// "function reverseRegistrar() view returns (address)",
func (c *RegistrarController) ReverseRegistrar() (common.Address, error) {
	return c.Contract.ReverseRegistrar(nil)
}

// "function supportsInterface(bytes4) pure returns (bool)",
func (c *RegistrarController) SupportsInterface(opts *bind.CallOpts, interfaceID [4]byte) (bool, error) {
	return c.Contract.SupportsInterface(nil, interfaceID)
}

// "function transferOwnership(address)",
func (c *RegistrarController) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return c.Contract.TransferOwnership(opts, newOwner)
}

// IsValid returns true if the domain is considered valid by the controller.
func (c *RegistrarController) IsValid(domain string) (bool, error) {
	name, err := UnqualifiedName(domain, c.domain)
	if err != nil {
		return false, fmt.Errorf("invalid name %s", domain)
	}
	return c.Contract.Valid(nil, name)
}

// "function withdraw()"
func (c *RegistrarController) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return c.Contract.Withdraw(opts)
}

// Reveal reveals a commitment to register a domain.
func (c *RegistrarController) Reveal(opts *bind.TransactOpts, domain string, owner common.Address, duration *big.Int, secret [32]byte) (*types.Transaction, error) {
	name, err := UnqualifiedName(domain, c.domain)
	if err != nil {
		return nil, fmt.Errorf("invalid name %s", domain)
	}

	if opts == nil {
		return nil, errors.New("transaction options required")
	}
	if opts.Value == nil {
		return nil, errors.New("no ether supplied with transaction")
	}

	commitTS, err := c.CommitmentTime(name, owner, duration, secret)
	if err != nil {
		return nil, err
	}
	if commitTS.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("no commitment present")
	}
	commit := time.Unix(commitTS.Int64(), 0)

	minCommitIntervalTS, err := c.MinCommitmentInterval()
	if err != nil {
		return nil, err
	}
	minCommitInterval := time.Duration(minCommitIntervalTS.Int64()) * time.Second
	minRevealTime := commit.Add(minCommitInterval)
	if time.Now().Before(minRevealTime) {
		return nil, errors.New("commitment too young to reveal")
	}

	maxCommitIntervalTS, err := c.MaxCommitmentInterval()
	if err != nil {
		return nil, err
	}
	maxCommitInterval := time.Duration(maxCommitIntervalTS.Int64()) * time.Second
	maxRevealTime := commit.Add(maxCommitInterval)
	if time.Now().After(maxRevealTime) {
		return nil, errors.New("commitment too old to reveal")
	}

	// Calculate the duration given the rent cost and the value
	// costPerSecond, err := c.RentCost(domain)
	// if err != nil {
	// 	return nil, errors.New("failed to obtain rent cost")
	// }
	// // duration := new(big.Int).Div(opts.Value, costPerSecond)

	// Ensure duration is greater than minimum duration
	minDuration, err := c.MinRegistrationDuration()
	if err != nil {
		return nil, err
	}
	if big.NewInt(int64(minDuration.Seconds())).Cmp(duration) >= 0 {
		return nil, fmt.Errorf("not enough funds to cover minimum duration of %v", minDuration)
	}

	config := getConfig()
	return c.Contract.Register(opts, name, owner, duration, secret, config.commitmentData.publicResover, config.commitmentData.calldata, config.commitmentData.reverseRecord, config.commitmentData.fuses, config.commitmentData.wrapperExpiry)
}

// Renew renews a registered domain.
func (c *RegistrarController) Renew(opts *bind.TransactOpts, domain string) (*types.Transaction, error) {
	name, err := UnqualifiedName(domain, c.domain)
	if err != nil {
		return nil, fmt.Errorf("invalid name %s", domain)
	}

	// See if we're registered at all - fetch the owner to find out
	registry, err := NewRegistry(c.backend)
	if err != nil {
		return nil, err
	}
	owner, err := registry.Owner(domain)
	if err != nil {
		return nil, err
	}
	if owner == UnknownAddress {
		return nil, fmt.Errorf("%s not registered", domain)
	}

	// Calculate the duration given the rent cost and the value
	costPerSecond, err := c.RentCost(domain)
	if err != nil {
		return nil, errors.New("failed to obtain rent cost")
	}
	duration := new(big.Int).Div(opts.Value, costPerSecond)

	return c.Contract.Renew(opts, name, duration)
}
