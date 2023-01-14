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
	//Check that the ENSRegistry has set the Deployer Address as the owner
	deployerAddress := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	nameHash, err := NameHash("resolver")
	// assert.Equal(t, err, nil, "Error getting country Namehash")
	// labelHash, err := LabelHash("resolver")
	// assert.Equal(t, err, nil, "Error getting resolver Labelhash")
	// sha := sha3.NewLegacyKeccak256()
	// // nodeHash := sha.Sum(append(nameHash[:], labelHash[:]))
	// nodeHash := sha3.NewLegacyKeccak256((append(nameHash[:], labelHash[:])))
	// // resolverNode := "0x84a0d1d9a4103b566373dcaf0bd2d482e763a581452c91a0eafb77bb49a1a71d")
	ensRegistry, err := ensregistry.NewContract(config.ENSRegistry, config.client)
	assert.Equal(t, err, nil, "Error getting ENSRegistry")
	resolverOwner, err := ensRegistry.Owner(nil, nameHash)
	assert.Equal(t, err, nil, "Error getting resolver node owner form ENSRegistry")
	assert.Equal(t, resolverOwner, deployerAddress, "Incorrect Owner for resolver node")

}
