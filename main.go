package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/hekmon/transmissionrpc"
	"github.com/mitchellh/go-homedir"
)

const configfilename = "~/.config/tmclient.json"

type TorrentHost struct {
	Hostname     string
	Username     string
	Password     string
	HTTPS        bool
	Port         uint16
	DownloadPath string
	FinalPath    string
}

type HostConfig map[string]TorrentHost

var (
	Hosts = make(HostConfig)
)

func init() {
	filename, _ := homedir.Expand(configfilename)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Unable to open config file %s: %s", filename, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Hosts)
	if err != nil {
		log.Fatalf("Unable to parse config file %s: %s", filename, err)
	}
}

func (h *TorrentHost) Connect() (*transmissionrpc.Client, error) {
	client, err := transmissionrpc.New(h.Hostname, h.Username, h.Password,
		&transmissionrpc.AdvancedConfig{
			HTTPS: h.HTTPS,
			Port:  h.Port,
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
	return client, nil
}

func main() {
	host := Hosts["two"]
	transmissionbt, err := host.Connect()
	if err != nil {
		log.Fatal(err)
	}
	torrents, err := transmissionbt.TorrentGetAll()
	for _, torrent := range torrents {
		fmt.Println()
		spew.Dump(torrent.Name, torrent.DoneDate, torrent.LeftUntilDone, torrent.PercentDone)
	}
}
