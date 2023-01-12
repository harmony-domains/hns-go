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

package onens

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

// Name represents an ENS name, for example 'foo.bar.eth'.
type config struct {
	Registrant           common.Address
	Controller           common.Address
	Resolver             common.Address
	Expiry               time.Time
	RegistrationInterval time.Duration
	clientURL            string
	client               *ethclient.Client
}

func getConfig() *config {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	config := &config{}
	config.Registrant = common.HexToAddress(viper.GetString("DSREGISTRANT"))
	config.Controller = common.HexToAddress(viper.GetString("DSREGISTRANT"))
	config.Resolver = common.HexToAddress(viper.GetString("DSRESOLVER"))
	config.Expiry = time.Unix(viper.GetInt64("DSREGISTRATION_INTERVAL"), 0)
	config.RegistrationInterval = viper.GetDuration("DSREGISTRATION_INTERVAL") * time.Second
	config.clientURL = viper.GetString("CLIENT_URL")
	config.client, _ = ethclient.Dial(config.clientURL)
	return config
}
