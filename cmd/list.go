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
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/hekmon/transmissionrpc"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func PrintTorrent(t *transmissionrpc.Torrent, width uint) {
	var (
		state aurora.Value
	)

	switch *t.Status {
	case transmissionrpc.TorrentStatusStopped:
		// 🛑 ❌
		if *t.PercentDone < 1 {
			state = aurora.BrightRed("✗")
		} else {
			state = aurora.Green("✗")
		}
	case transmissionrpc.TorrentStatusCheckWait:
		state = aurora.Red("⏳")
	case transmissionrpc.TorrentStatusCheck:
		state = aurora.BgCyan("🔍")
	case transmissionrpc.TorrentStatusDownloadWait:
		state = aurora.Green("↯")
	case transmissionrpc.TorrentStatusDownload:
		state = aurora.Green("↯")
	case transmissionrpc.TorrentStatusSeedWait:
		state = aurora.Green("🛏")
	case transmissionrpc.TorrentStatusSeed:
		state = aurora.Green("✓")
	case transmissionrpc.TorrentStatusIsolated:
		state = aurora.Magenta("❓")
	}

	fmt.Printf("%s %3.1f%% %s\n", state, *t.PercentDone*100, *t.Name)
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [host]",
	Short: "List torrents on host",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("No host specified")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := Connect(viper.GetViper(), args[0])
		if err != nil {
			log.Fatal(err)
		}
		torrents, err := client.TorrentGetAll()
		if err != nil {
			log.Fatal(err)
		}
		sort.SliceStable(torrents, func(i int, j int) bool {
			switch {
			case *torrents[i].IsFinished && *torrents[j].IsFinished:
				return torrents[i].AddedDate.Before(*torrents[j].AddedDate)
			case *torrents[i].PercentDone < 1 && *torrents[j].PercentDone == 1:
				return true
			case *torrents[i].PercentDone == 1 && *torrents[j].PercentDone < 1:
				return false
			default:
				return *torrents[i].PercentDone < *torrents[j].PercentDone
			}
		})
		for _, t := range torrents {
			PrintTorrent(t, 80) // TODO: use real terminal width here
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// ▱▰
