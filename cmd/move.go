// Copyright © 2019 Heiko Reese <mail@heiko-reese.de>
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

package cmd

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// moveCmd represents the move command
var moveCmd = &cobra.Command{
	Use:   "move [host]",
	Short: "Move downloaded torrents to preconfigured directory",
	Long:  `If host is not specified, all configured hosts are checked`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			knownHosts = make(map[string]bool)
			hosts      []string
		)
		for _, h := range GetHosts(viper.GetViper()) {
			knownHosts[h] = true
		}
		for _, h := range args {
			_, exists := knownHosts[h]
			if !exists {
				_, _ = fmt.Fprintf(os.Stderr, "Host %s unknown", h)
				os.Exit(2)
			}
		}
		if len(args) == 0 {
			hosts = GetHosts(viper.GetViper())
		} else {
			hosts = args
		}

		for _, host := range hosts {
			keydownload := fmt.Sprintf("hosts.%s.downloadpath", host)
			keyfinal := fmt.Sprintf("hosts.%s.finalpath", host)
			if viper.IsSet(keydownload) && viper.IsSet(keydownload) {
				downloadpath := viper.GetString(keydownload)
				finalpath := viper.GetString(keyfinal)
				client, err := Connect(viper.GetViper(), host)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to %s (%s), skipping\n", host, err)
					continue
				}
				torrents, err := client.TorrentGetAll()
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Unable to get list of torrents for %s (%s), skipping\n", host, err)
					continue
				}
				for _, torrent := range torrents {
					if *torrent.DownloadDir == downloadpath && *torrent.PercentDone >= 1 {
						err = client.TorrentSetLocation(*torrent.ID, finalpath, true)
						if err != nil {
							_, _ = fmt.Fprintf(os.Stderr, "Unable to move %s: %s\n", *torrent.Name, aurora.Bold(err))
						} else {
							_, _ = fmt.Printf("[%s] Moved %s\n", aurora.Italic(host), aurora.Bold(*torrent.Name))
						}
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(moveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// moveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// moveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
