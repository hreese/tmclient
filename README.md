# About

tmclient is a limited client for the [transmission](https://transmissionbt.com/)
server. It's focused on a few tasks that enable support handling lots if bittorrent
files and some automation.

# Configuration

Create a file `.tmclient.json` in your `$HOME` (that's usually `C:\Users\yourname`
on Windows) which declares all transmission hosts:

```
{
  "hosts": {
    "server_one": {
      "Hostname": "first-host.example.org",
      "Username": "rpc-username",
      "Password": "rpc-password",
      "HTTPS": true,
      "Port": 9091,
      "DownloadPath": "/home/transmission/incoming",
      "FinalPath": "/home/transmission/finished"
    },
    "number_two": {
      "Hostname": "another-host.example.com",
      "Username": "rpc-username",
      "Password": "rpc-password",
      "HTTPS": false,
      "Port": 9091
    }
  }
}

```

You may skip `DownloadPath` amd `FinalPath`. Make sure that the host's names
(`server_one` and `number_two` in this example) are unique.

# Usage

Calling `tmclient` (Windows user may need to add `.exe` to all calls) will list
all configured hosts. All commands and subcommands are documented online (use
`help` or `--help`).

To add all torrent files in the current directory to host number_two simply call
`tmclient add number_two *.torrent`. Files that are successfully added are removed
unless the `--keep` flag is used.

`tmclient list number_two ` produces a rudimentary listing of all active torrents.

`tmclient move [hosts]` will move all completed torrents from hosts (or all hosts
when no host is specified) from `DownloadPath` to `FinalPath`. This is useful if
there is an automated copy mechanism in place (think `while (sleep 60); do rsync …; done`)
to copy files to their final destination. Just run `tmclient move` periodically
in a cronjob/systemd timer to get files as soon as they finish downloading.