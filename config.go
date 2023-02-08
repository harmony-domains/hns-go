// Copyright John Whitton https://github.com/john_whitton
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
Addresses can be obtained from 1ns-deployer by running when deploying the contracts
the comand used is `yarn deploy --network local` replacing `local` with the network you are deploying on

Below is sample output from a local deploy

oracleDeployer: 0x5FbDB2315678afecb367f032d93F642f64180aa3
- priceOracle: 0xB7A5bd0345EF1Cc5E66bf61BdeC17D2461fBd968
- usdOracle: 0xa16E02E87b7454126E5E10d957A927A7F5B5d2be
deploying "ENSDeployer" (tx: 0x581ad4a55efc30a9c2530876e46ece8b903f6c93ef0f8777587b7c2436bc767c)...: deployed at 0x0165878A594ca255338adfa4d48449f69242Eb8F with 18676739 gas
deployer account 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
ENSDeployer deployed to: 0x0165878A594ca255338adfa4d48449f69242Eb8F
- ens deployed to: 0x3B02fF1e626Ed7a8fd6eC5299e2C54e1421B626B
- fifsRegistrar deployed to: 0xBA12646CC07ADBe43F8bD25D83FB628D29C8A762
- reverseRegistrar deployed to: 0x7ab4C4804197531f7ed6A6bc0f0781f706ff7953
- baseRegistrar deployed to: 0xc8CB5439c767A63aca1c01862252B2F3495fDcFE
- metadataService deployed to: 0xD79aE87F2c003Ec925fB7e9C11585709bfe41473
- nameWrapper deployed to: 0xB7aa4c318000BB9bD16108F81C40D02E48af1C42
- registrarController deployed to: 0x12653A08808F651D5BB78514F377d3BD5E17934C
- publicResolver deployed to: 0xCaA29B65446aBF1A513A178402A0408eB3AEee75
- universalResolver deployed to: 0x09F428b7D940ED8Bff862e81a103bf022F5E50F0
tx 0x20a148fac52a922e4956ec21330dcc1e39307d0734dd23cc301e68438cdbdba9
*/

package onens

import (
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

// Constants
var zeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")

// Configuration
type testAccounts struct {
	deployerAddress     common.Address
	deployerPrivateKey  *ecdsa.PrivateKey
	operatorAAddress    common.Address
	operatorAPrivateKey *ecdsa.PrivateKey
	operatorBAddress    common.Address
	operatorBPrivateKey *ecdsa.PrivateKey
	operatorCAddress    common.Address
	operatorCPrivateKey *ecdsa.PrivateKey
	aliceAddress        common.Address
	alicePrivateKey     *ecdsa.PrivateKey
	bobAddress          common.Address
	bobPrivateKey       *ecdsa.PrivateKey
	carolAddress        common.Address
	carolPrivateKey     *ecdsa.PrivateKey
	doraAddress         common.Address
	doraPrivateKey      *ecdsa.PrivateKey
	ernieAddress        common.Address
	erniePrivateKey     *ecdsa.PrivateKey
	fredAddress         common.Address
	fredPrivateKey      *ecdsa.PrivateKey
}

// Test configuration Structure
type tconfigStruct struct {
	testAccounts
	PriceOracle          common.Address
	USDOracle            common.Address
	Registry             common.Address
	FIFSRegistrar        common.Address
	ReverseRegistrar     common.Address
	BaseRegistrar        common.Address
	MetadataService      common.Address
	NameWrapper          common.Address
	RegistrarController  common.Address
	PublicResolver       common.Address
	UniversalResolver    common.Address
	Registrant           common.Address
	Expiry               time.Time
	RegistrationInterval time.Duration
	clientURL            string
	chainID              int64
	client               *ethclient.Client
	TLD                  string
	duration             *big.Int
}

// onens Default commitment data used for registration comittment data
type commitmentData struct {
	secret        [32]byte
	publicResover common.Address
	calldata      [][]byte
	reverseRecord bool
	fuses         uint32
	wrapperExpiry uint64
}

// Configuration Structure
type configStruct struct {
	Registry common.Address
	commitmentData
}

var tconfig *tconfigStruct = getTConfig()
var tclient *ethclient.Client = tconfig.client

var config *configStruct = getConfig()

// Get Test Configuration
func getTConfig() *tconfigStruct {
	tconfig := &tconfigStruct{}
	// Read test config from environment file
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	// set test accounts
	tconfig.testAccounts.deployerAddress = common.HexToAddress(viper.GetString("TEST_ADDRESS_DEPLOYER"))
	tconfig.testAccounts.deployerPrivateKey, _ = crypto.HexToECDSA(viper.GetString("TEST_PRIVATE_KEY_DEPLOYER"))
	tconfig.testAccounts.operatorAAddress = common.HexToAddress(viper.GetString("TEST_ADDRESS_OPEATORA"))
	tconfig.testAccounts.operatorAPrivateKey, _ = crypto.HexToECDSA(viper.GetString("TEST_PRIVATE_KEY_OPERATORA"))
	tconfig.testAccounts.operatorBAddress = common.HexToAddress(viper.GetString("TEST_ADDRESS_OPERATORB"))
	tconfig.testAccounts.operatorBPrivateKey, _ = crypto.HexToECDSA(viper.GetString("TEST_PRIVATE_KEY_OPERATORB"))
	tconfig.testAccounts.operatorCAddress = common.HexToAddress(viper.GetString("TEST_ADDRESS_OPERATORC"))
	tconfig.testAccounts.operatorCPrivateKey, _ = crypto.HexToECDSA(viper.GetString("TEST_PRIVATE_KEY_OPERATORC"))
	tconfig.testAccounts.aliceAddress = common.HexToAddress(viper.GetString("TEST_ADDRESS_ALICE"))
	tconfig.testAccounts.alicePrivateKey, _ = crypto.HexToECDSA(viper.GetString("TEST_PRIVATE_KEY_ALICE"))
	tconfig.testAccounts.bobAddress = common.HexToAddress(viper.GetString("TEST_ADDRESS_BOB"))
	tconfig.testAccounts.bobPrivateKey, _ = crypto.HexToECDSA(viper.GetString("TEST_PRIVATE_KEY_BOB"))
	tconfig.testAccounts.carolAddress = common.HexToAddress(viper.GetString("TEST_ADDRESS_CAROL"))
	tconfig.testAccounts.carolPrivateKey, _ = crypto.HexToECDSA(viper.GetString("TEST_PRIVATE_KEY_CAROL"))
	tconfig.testAccounts.doraAddress = common.HexToAddress(viper.GetString("TEST_ADDRESS_DORA"))
	tconfig.testAccounts.doraPrivateKey, _ = crypto.HexToECDSA(viper.GetString("TEST_PRIVATE_KEY_DORA"))
	tconfig.testAccounts.ernieAddress = common.HexToAddress(viper.GetString("TEST_ADDRESS_ERNIE"))
	tconfig.testAccounts.erniePrivateKey, _ = crypto.HexToECDSA(viper.GetString("TEST_PRIVATE_KEY_FRANK"))
	tconfig.testAccounts.fredAddress = common.HexToAddress(viper.GetString("TEST_ADDRESS_FRED"))
	tconfig.testAccounts.fredPrivateKey, _ = crypto.HexToECDSA(viper.GetString("TEST_PRIVATE_KEY_ERNIE"))
	// set additional test configuration
	tconfig.PriceOracle = common.HexToAddress(viper.GetString("TEST_PRICE_ORACLE"))
	tconfig.USDOracle = common.HexToAddress(viper.GetString("TEST_USD_ORACLE"))
	tconfig.Registry = common.HexToAddress(viper.GetString("TEST_ENS_REGISTRY"))
	tconfig.FIFSRegistrar = common.HexToAddress(viper.GetString("TEST_FIFS_REGISTRAR"))
	tconfig.ReverseRegistrar = common.HexToAddress(viper.GetString("TEST_REVERSE_REGISTRAR"))
	tconfig.BaseRegistrar = common.HexToAddress(viper.GetString("TEST_BASE_REGISTRAR"))
	tconfig.MetadataService = common.HexToAddress(viper.GetString("TEST_METADATA_SERVICE"))
	tconfig.NameWrapper = common.HexToAddress(viper.GetString("TEST_NAME_WRAPPER"))
	tconfig.RegistrarController = common.HexToAddress(viper.GetString("TEST_REGISTRAR_CONTROLLER"))
	tconfig.PublicResolver = common.HexToAddress(viper.GetString("TEST_PUBLIC_RESOLVER"))
	tconfig.UniversalResolver = common.HexToAddress(viper.GetString("TEST_UNIVERSAL_RESOLVER"))
	tconfig.Expiry = time.Unix(viper.GetInt64("TEST_EXPIRY"), 0)
	tconfig.RegistrationInterval = viper.GetDuration("TEST_REGISTRATION_INTERVAL") * time.Second
	tconfig.clientURL = viper.GetString("TEST_CLIENT_URL")
	tconfig.chainID = viper.GetInt64("TEST_CHAIN_ID")
	client, err := ethclient.Dial(tconfig.clientURL)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to connect to Ethereum Client")
		fmt.Println(tconfig.clientURL)
		os.Exit(1)
		log.Fatal(err)
	}
	tconfig.client = client
	tconfig.TLD = viper.GetString("TLD")
	tconfig.duration = big.NewInt(viper.GetInt64("TEST_DURATION"))
	return tconfig
}

// Get onens configuration
func getConfig() *configStruct {
	config := &configStruct{}
	// Read test config from environment file
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	// set configuration
	config.Registry = common.HexToAddress(viper.GetString("ENS_REGISTRY"))
	// Set Commitment Data
	config.commitmentData.secret = [32]byte{}
	config.commitmentData.publicResover = common.HexToAddress(viper.GetString("PUBLIC_RESOLVER"))
	config.commitmentData.calldata = [][]byte{}
	config.commitmentData.reverseRecord = false
	config.commitmentData.fuses = 0
	config.commitmentData.wrapperExpiry = math.MaxUint64
	return config
}

func decode(pemEncoded string) *ecdsa.PrivateKey {
	// block, _ := pem.Decode([]byte(pemEncoded))
	// x509Encoded := block.Bytes
	x509Encoded := []byte(pemEncoded)

	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

	return privateKey
}
