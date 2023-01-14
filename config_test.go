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
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jw-1ns/go-1ns/contracts/baseregistrar"
	"github.com/jw-1ns/go-1ns/contracts/ensregistry"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	config := getConfig()
	client = config.client
	// Test we can connect to the client
	// Test we can get configured contract addresses and they match
	//Functional Tests

	// Check Base Registrar is pointing to the ENS Registry
	baseRegistrar, err := baseregistrar.NewContract(config.BaseRegistrar, config.client)
	assert.Equal(t, err, nil, "Error getting BaseRegistrar")
	baseRegistrarENS, err := baseRegistrar.Ens(nil)
	// baseRegistrarAddress, err := baseRegistrar.Address()
	fmt.Printf("baseRegistrar.Ens: %v\n", baseRegistrarENS)
	// fmt.Printf("baseRegistrar: %v\n", &baseRegistrar)
	assert.Equal(t, err, nil, "Error getting ENS from baseRegistrar")
	assert.Equal(t, baseRegistrarENS, config.ENSRegistry, "Incorrect ENS for baseRegistrar")

	//Check that the ENSRegistry has test owners for domains set correctly
	// await registerDomain('test', alice, await ensDeployer.publicResolver(), await ensDeployer.registrarController())
	// await registerDomain('resolver', bob, await ensDeployer.publicResolver(), await ensDeployer.registrarController())
	// deployerAddress := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	aliceAddress := common.HexToAddress("0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65")
	bobAddress := common.HexToAddress("0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc")
	countryNameHash, err := NameHash("country")
	assert.Equal(t, err, nil, "Error getting Namehash for country")
	testNameHash, err := NameHash("test")
	assert.Equal(t, err, nil, "Error getting Namehash for test")
	resolverNameHash, err := NameHash("resolver")
	assert.Equal(t, err, nil, "Error getting Namehash for resolver")
	ensRegistry, err := ensregistry.NewContract(config.ENSRegistry, config.client)
	assert.Equal(t, err, nil, "Error getting ENSRegistry")
	// country is owned by the BaseRegistrar
	countryOwner, err := ensRegistry.Owner(nil, countryNameHash)
	assert.Equal(t, err, nil, "Error getting resolver node owner form ENSRegistry")
	assert.Equal(t, countryOwner, config.BaseRegistrar, "Incorrect Owner for country node")
	// test is owned by Alice
	testOwner, err := ensRegistry.Owner(nil, testNameHash)
	assert.Equal(t, err, nil, "Error getting resolver node owner form ENSRegistry")
	assert.Equal(t, testOwner, aliceAddress, "Incorrect Owner for test node")
	// resolver is owned by Bob
	resolverOwner, err := ensRegistry.Owner(nil, resolverNameHash)
	assert.Equal(t, err, nil, "Error getting resolver node owner form ENSRegistry")
	assert.Equal(t, resolverOwner, bobAddress, "Incorrect Owner for resolver node")

}

// assert.Equal(t, err, nil, "Error getting country Namehash")
// labelHash, err := LabelHash("resolver")
// assert.Equal(t, err, nil, "Error getting resolver Labelhash")
// sha := sha3.NewLegacyKeccak256()
// // nodeHash := sha.Sum(append(nameHash[:], labelHash[:]))
// nodeHash := sha3.NewLegacyKeccak256((append(nameHash[:], labelHash[:])))
// // resolverNode := "0x84a0d1d9a4103b566373dcaf0bd2d482e763a581452c91a0eafb77bb49a1a71d")
