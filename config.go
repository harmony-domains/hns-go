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

package hns

import (
	"time"

	"github.com/ava-labs/coreth/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
)

// Name represents an ENS name, for example 'foo.bar.eth'.
type config struct {
	dsRegistrant           common.Address
	dsController           common.Address
	dsResolver             common.Address
	dsExpiry               time.Time
	dsRegistrationInterval time.Duration
	client                 ethclient.Client
}

func getConfig() *config {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	config := &config{}
	config.dsRegistrant = common.HexToAddress(viper.GetString("DSREGISTRANT"))
	config.dsController = common.HexToAddress(viper.GetString("DSREGISTRANT"))
	config.dsResolver = common.HexToAddress(viper.GetString("DSRESOLVER"))
	config.dsExpiry = time.Unix(viper.GetInt64("DSREGISTRATION_INTERVAL"), 0)
	config.dsRegistrationInterval = viper.GetDuration("DSREGISTRATION_INTERVAL") * time.Second
	config.client, _ = ethclient.Dial(viper.GetString("CLIENT_URL"))

	// 	dsController := common.HexToAddress("a303ddc620aa7d1390baccc8a495508b183fab59")
	// 	dsResolver := common.HexToAddress("DaaF96c344f63131acadD0Ea35170E7892d3dfBA")
	// 	dsExpiry := time.Unix(4741286688, 0)
	// 	dsRegistrationInterval := 60 * time.Second

	// 	client, _ := ethclient.Dial("https://ropsten.infura.io/v3/831a5442dc2e4536a9f8dee4ea1707a6")
	return config
}
