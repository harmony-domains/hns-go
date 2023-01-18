// Copyright 2019, 2022 Weald Technology Trading
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
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jw-1ns/go-1ns/contracts/baseregistrar"
	"github.com/jw-1ns/go-1ns/contracts/ensregistry"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	// config := getConfig()
	// Test we can connect to the client
	// Test we can get configured contract addresses and they match
	//Functional Tests

	// Check Base Registrar is pointing to the ENS Registry
	baseRegistrar, err := baseregistrar.NewContract(tconfig.BaseRegistrar, tclient)
	assert.Equal(t, err, nil, "Error getting BaseRegistrar")
	baseRegistrarENS, err := baseRegistrar.Ens(nil)
	assert.Equal(t, err, nil, "Error getting ENS from baseRegistrar")
	assert.Equal(t, baseRegistrarENS, tconfig.ENSRegistry, "Incorrect ENS for baseRegistrar")

	//Check that the ENSRegistry has test owners for domains set correctly
	deployerAddress := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	ensRegistry, err := ensregistry.NewContract(tconfig.ENSRegistry, tclient)
	assert.Equal(t, err, nil, "Error getting ENSRegistry")
	// country is owned by the BaseRegistrar
	countryNameHash, err := NameHash("country")
	assert.Equal(t, err, nil, "Error getting Namehash for country")
	countryOwner, err := ensRegistry.Owner(nil, countryNameHash)
	assert.Equal(t, err, nil, "Error getting resolver node owner from ENSRegistry")
	assert.Equal(t, countryOwner, tconfig.BaseRegistrar, "Incorrect Owner for country node")
	// resolver is owned by Deployer
	resolverNameHash, err := NameHash("resolver")
	assert.Equal(t, err, nil, "Error getting Namehash for resolver")
	resolverOwner, err := ensRegistry.Owner(nil, resolverNameHash)
	assert.Equal(t, err, nil, "Error getting resolver node owner from ENSRegistry")
	assert.Equal(t, resolverOwner, deployerAddress, "Incorrect Owner for resolver node")
	// test.country is owned by NameWrapper
	tcNameHash, err := NameHash("test.country")
	assert.Equal(t, err, nil, "Error getting nameHash for test.country")
	testCountryOwner, err := ensRegistry.Owner(nil, tcNameHash)
	assert.Equal(t, err, nil, "Error getting test.country node owner fromm ENSRegistry")
	assert.Equal(t, testCountryOwner, tconfig.NameWrapper, "Incorrect Owner for test.country node")
	// test is owned by Namewrapper using BaseRegistrar TokenId
	testLabelHash, err := LabelHash("test")
	assert.Equal(t, err, nil, "Error getting LabelHash for test")
	testLabelHashBigInt := new(big.Int).SetBytes(testLabelHash[:])
	testOwnerBaseRegistrar, err := baseRegistrar.OwnerOf(nil, testLabelHashBigInt)
	assert.Equal(t, err, nil, "Error getting test owner from ENSRegistry")
	assert.Equal(t, testOwnerBaseRegistrar, tconfig.NameWrapper, "Incorrect Owner for test.country node")
	// testxyz is owned by Namewrapper using BaseRegistrar TokenId
	testxyzLabelHash, err := LabelHash("testxyz")
	assert.Equal(t, err, nil, "Error getting LabelHash for testxyz")
	testxyzLabelHashBigInt := new(big.Int).SetBytes(testxyzLabelHash[:])
	testxyzOwnerBaseRegistrar, err := baseRegistrar.OwnerOf(nil, testxyzLabelHashBigInt)
	assert.Equal(t, err, nil, "Error getting testxyz owner from ENSRegistry")
	assert.Equal(t, testxyzOwnerBaseRegistrar, zeroAddress, "Incorrect Owner for testxyz.country node")

	// unregistered Tier 2 have no owners
	unregisteredNameHash, err := NameHash("unregistered.country")
	assert.Equal(t, err, nil, "Error getting nameHash for unregistered.country")
	unregisteredCountryOwner, err := ensRegistry.Owner(nil, unregisteredNameHash)
	assert.Equal(t, err, nil, "Error getting unregistered.country node owner from ENSRegistry")
	assert.Equal(t, unregisteredCountryOwner, zeroAddress, "Incorrect Owner for unregistered.country node")

}
