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
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveEmpty(t *testing.T) {
	_, err := Resolve(tclient, "")
	assert.NotNil(t, err, "Resolved empty name")
}

func TestResolveZero(t *testing.T) {
	_, err := Resolve(tclient, "0")
	assert.NotNil(t, err, "Resolved empty name")
}

func TestResolveNotPresent(t *testing.T) {
	_, err := Resolve(tclient, "sirnotappearinginthisregistry.country")
	require.NotNil(t, err, "Resolved name that does not exist")
	assert.Equal(t, "unregistered name", err.Error(), "Unexpected error")
}

// func TestResolveNoResolver(t *testing.T) {
// 	_, err := Resolve(tclient, "noresolver.country")
// 	require.NotNil(t, err, "Resolved name without a resolver")
// 	assert.Equal(t, "no resolver", err.Error(), "Unexpected error")
// }

// func TestResolveBadResolver(t *testing.T) {
// 	_, err := Resolve(tclient, "resolvestozero.country")
// 	require.NotNil(t, err, "Resolved name with a bad resolver")
// 	assert.Equal(t, "no address", err.Error(), "Unexpected error")
// }

func TestResolveTestCountry(t *testing.T) {
	_, err := Resolve(tclient, "test.country")
	assert.Equal(t, err.Error(), "no address")
}

// func TestResolveTestEnsTest(t *testing.T) {
// 	expected := "388ea662ef2c223ec0b047d41bf3c0f362142ad5"
// 	actual, err := Resolve(tclient, "test.enstest.country")
// 	require.Nil(t, err, "Error resolving name")
// 	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
// }

func TestResolveResolverCountry(t *testing.T) {
	expected := "caa29b65446abf1a513a178402a0408eb3aeee75"
	actual, err := Resolve(tclient, "resolver.country")
	require.Nil(t, err, "Error resolving name")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

func TestResolveCountry(t *testing.T) {
	expected := "000000000000000000000000000000000000000c"
	actual, err := Resolve(tclient, "country")
	require.Nil(t, err, "Error resolving name")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

// func TestResolveNickJo1nson(t *testing.T) {
// 	expected := "70abd981e01ad3e6eb1726a5001000877ab04417"
// 	actual, err := Resolve(tclient, "nickjo1nson.country")
// 	require.Nil(t, err, "Error resolving name")
// 	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
// }

func TestResolveAddress(t *testing.T) {
	expected := "5ffc014343cd971b7eb70732021e26c35b744cc4"
	actual, err := Resolve(tclient, "0x5ffc014343cd971b7eb70732021e26c35b744cc4")
	require.Nil(t, err, "Error resolving address")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

func TestResolveShortAddress(t *testing.T) {
	expected := "0000000000000000000000000000000000000001"
	actual, err := Resolve(tclient, "0x1")
	require.Nil(t, err, "Error resolving address")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

func TestResolveHexString(t *testing.T) {
	_, err := Resolve(tclient, "0xe32c6d1a964749b6de2130e20daed821a45b9e7261118801ff5319d0ffc6b54a")
	assert.NotNil(t, err, "Resolved too-long hex string")
}

// func TestReverseResolveTestEnsTest(t *testing.T) {
// 	expected := "1ns.country"
// 	address := common.HexToAddress("0x388ea662ef2c223ec0b047d41bf3c0f362142ad5")
// 	actual, err := ReverseResolve(client, address)
// 	require.Nil(t, err, "Error resotlving address")
// 	assert.Equal(t, expected, actual, "Did not receive expected result")
// }
