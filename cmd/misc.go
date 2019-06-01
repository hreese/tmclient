/*
Copyright 2019 Heiko Reese

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"errors"
	"fmt"
	"github.com/hekmon/transmissionrpc"
	"github.com/spf13/viper"
)

/*type TorrentHost struct {
	Hostname     string
	Username     string
	Password     string
	HTTPS        bool
	Port         uint16
	DownloadPath string
	FinalPath    string
}*/

func Connect (v *viper.Viper, hostname string) (*transmissionrpc.Client, error) {
	if !v.InConfig("hosts") {
		return nil, errors.New("No hosts configured")
	}
	//for key, _ := range v.Sub("hosts").AllSettings()
	if !v.Sub("hosts").IsSet(hostname) {
		return nil, errors.New(fmt.Sprintf("No host %s configured", hostname))
	}
	hostviper := v.Sub("hosts").Sub(hostname)
	client, err := transmissionrpc.New(
		hostviper.GetString("hostname"),
		hostviper.GetString("username"),
		hostviper.GetString("password"),
		&transmissionrpc.AdvancedConfig{
			HTTPS: hostviper.GetBool("https"),
			Port: uint16(hostviper.GetUint32("port")),
		})
	if err != nil {
		return nil, err
	}
	ok, serverVersion, serverMinimumVersion, err := client.RPCVersion()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New(fmt.Sprintf("Remote transmission RPC version (v%d) is incompatible with the transmission library (v%d): remote needs at least v%d",
			serverVersion, transmissionrpc.RPCVersion, serverMinimumVersion))
	}
	return client, nil}


func GetHosts(v *viper.Viper) []string {
	var hostnames []string
	if v.InConfig("hosts") {
		for key, _ := range v.Sub("hosts").AllSettings() {
			hostnames = append(hostnames, key)
		}
	}
	return hostnames
}