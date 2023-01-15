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
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jw-1ns/go-1ns/contracts/baseregistrar"
	"github.com/jw-1ns/go-1ns/contracts/ensregistry"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	zeroAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")
	// config := getConfig()
	// Test we can connect to the client
	// Test we can get configured contract addresses and they match
	//Functional Tests

	// Check Base Registrar is pointing to the ENS Registry
	baseRegistrar, err := baseregistrar.NewContract(config.BaseRegistrar, config.client)
	assert.Equal(t, err, nil, "Error getting BaseRegistrar")
	baseRegistrarENS, err := baseRegistrar.Ens(nil)
	assert.Equal(t, err, nil, "Error getting ENS from baseRegistrar")
	assert.Equal(t, baseRegistrarENS, config.ENSRegistry, "Incorrect ENS for baseRegistrar")

	//Check that the ENSRegistry has test owners for domains set correctly
	deployerAddress := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	ensRegistry, err := ensregistry.NewContract(config.ENSRegistry, config.client)
	assert.Equal(t, err, nil, "Error getting ENSRegistry")
	// country is owned by the BaseRegistrar
	countryNameHash, err := NameHash("country")
	assert.Equal(t, err, nil, "Error getting Namehash for country")
	countryOwner, err := ensRegistry.Owner(nil, countryNameHash)
	assert.Equal(t, err, nil, "Error getting resolver node owner form ENSRegistry")
	assert.Equal(t, countryOwner, config.BaseRegistrar, "Incorrect Owner for country node")
	// resolver is owned by Deployer
	resolverNameHash, err := NameHash("resolver")
	assert.Equal(t, err, nil, "Error getting Namehash for resolver")
	resolverOwner, err := ensRegistry.Owner(nil, resolverNameHash)
	assert.Equal(t, err, nil, "Error getting resolver node owner form ENSRegistry")
	assert.Equal(t, resolverOwner, deployerAddress, "Incorrect Owner for resolver node")
	// test.country is owned by NameWrapper
	tcNameHash, err := nameHashCountry("test.country")
	assert.Equal(t, err, nil, "Error getting nameHashCountry for test.country")
	testCountryOwner, err := ensRegistry.Owner(nil, tcNameHash)
	assert.Equal(t, err, nil, "Error getting test.country node owner form ENSRegistry")
	assert.Equal(t, testCountryOwner, config.NameWrapper, "Incorrect Owner for test.country node")
	// unregistered Tier 2 have no owners
	unregisteredNameHash, err := nameHashCountry("unregistered.country")
	assert.Equal(t, err, nil, "Error getting nameHashCountry for unregistered.country")
	unregisteredCountryOwner, err := ensRegistry.Owner(nil, unregisteredNameHash)
	assert.Equal(t, err, nil, "Error getting unregistered.country node owner fromm ENSRegistry")
	assert.Equal(t, unregisteredCountryOwner, zeroAddress, "Incorrect Owner for unregistered.country node")

	// Additional Hardcoded Test
	// test.country get owener useing hardcoded values
	// test.country is owned by NameWrapper
	// testCountryNameHashHex := "6ccdbd41a174e9b5e34bffee7b0cc45c3ef17f8763cd491f14bc52dbb550b3b2"
	testCountryNameHash := []byte{108, 205, 189, 65, 161, 116, 233, 181, 227, 75, 255, 238, 123, 12, 196, 92, 62, 241, 127, 135, 99, 205, 73, 31, 20, 188, 82, 219, 181, 80, 179, 178}
	var testCountryNameHash32 [32]byte
	copy(testCountryNameHash32[:], testCountryNameHash)
	testCountryOwnerHardCoded, err := ensRegistry.Owner(nil, testCountryNameHash32)
	assert.Equal(t, err, nil, "Error getting test.country hardcoded node owner form ENSRegistry")
	assert.Equal(t, testCountryOwnerHardCoded, config.NameWrapper, "Incorrect Owner for test.country hardcoded node")

}
