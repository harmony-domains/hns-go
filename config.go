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
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

// Name represents an ENS name, for example 'foo.bar.country'.
type commitmentData struct {
	secret        [32]byte
	publicResover common.Address
	calldata      [][]byte
	reverseRecord bool
	fuses         uint32
	wrapperExpiry uint64
}
type config struct {
	PriceOracle          common.Address
	USDOracle            common.Address
	ENSRegistry          common.Address
	FIFSRegistrar        common.Address
	ReverseRegistrar     common.Address
	BaseRegistrar        common.Address
	MetadataService      common.Address
	NameWrapper          common.Address
	RegistrarController  common.Address
	PublicResover        common.Address
	UniversalResolver    common.Address
	Registrant           common.Address
	Controller           common.Address
	Resolver             common.Address
	Expiry               time.Time
	RegistrationInterval time.Duration
	clientURL            string
	client               *ethclient.Client
	TLD                  string
	commitmentData
}

func getConfig() *config {
	config := &config{}
	// Set Commitment Data
	config.commitmentData.secret = [32]byte{}
	config.commitmentData.publicResover = common.HexToAddress(viper.GetString("PUBLIC_RESOLVER"))
	config.commitmentData.calldata = [][]byte{}
	config.commitmentData.reverseRecord = false
	config.commitmentData.fuses = 0
	config.commitmentData.wrapperExpiry = math.MaxUint64

	// Read config from environment file
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	config.PriceOracle = common.HexToAddress(viper.GetString("PRICE_ORACLE"))
	config.USDOracle = common.HexToAddress(viper.GetString("USD_ORACLE"))
	config.ENSRegistry = common.HexToAddress(viper.GetString("ENS_REGISTRY"))
	config.FIFSRegistrar = common.HexToAddress(viper.GetString("FIFS_REGISTRAR"))
	config.ReverseRegistrar = common.HexToAddress(viper.GetString("REVERSE_REGISTRAR"))
	config.BaseRegistrar = common.HexToAddress(viper.GetString("BASE_REGISTRAR"))
	config.MetadataService = common.HexToAddress(viper.GetString("METADATA_SERVICE"))
	config.NameWrapper = common.HexToAddress(viper.GetString("NAME_WRAPPER"))
	config.RegistrarController = common.HexToAddress(viper.GetString("REGISTRAR_CONTROLLER"))
	config.PublicResover = common.HexToAddress(viper.GetString("PUBLIC_RESOLVER"))
	config.UniversalResolver = common.HexToAddress(viper.GetString("UNIVERSAL_RESOLVER"))
	config.RegistrationInterval = viper.GetDuration("DSREGISTRATION_INTERVAL") * time.Second
	config.clientURL = viper.GetString("CLIENT_URL")
	client, err := ethclient.Dial(config.clientURL)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to connect to Ethereum Client")
		fmt.Println(config.clientURL)
		os.Exit(1)
		log.Fatal(err)
	}
	config.client = client
	config.TLD = viper.GetString("TLD")
	return config
}
