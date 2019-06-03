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
	"log"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [host] [torrent…]",
	Short: "Add torrent files to host",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("No host specified")
		}
		if len(args) < 2 {
			return errors.New("No .torrent files to add")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := Connect(viper.GetViper(), args[0])
		if err != nil {
			log.Fatal(err)
		}
		for _, filename := range args[1:] {
			_, _ = fmt.Fprintf(os.Stderr, "Adding %s ", aurora.Bold(filename))
			_, err = client.TorrentAddFile(filename)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", aurora.Red("❌"), aurora.Bold(err))
				continue
			}
			_, _ = fmt.Fprintf(os.Stderr, "%s", aurora.Green("✔"))
			if viper.GetBool("keep") == false {
				err = os.Remove(filename)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, " %s: %s\n", aurora.Red("♻"), aurora.Bold(err))
				} else {
					_, _ = fmt.Fprintf(os.Stderr, " %s\n", aurora.Green("♻"))
				}
			} else {
				_, _ = fmt.Fprintln(os.Stderr)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("keep", "k", false, "Keep .torrent files after successful add (default is delete)")
	_ = viper.BindPFlag("keep", addCmd.Flags().Lookup("keep"))
}
