[![Build Status](https://img.shields.io/travis/LogicalOverflow/music-sync/master.svg?style=flat-square)](https://travis-ci.org/LogicalOverflow/music-sync)
[![Coverage Status](https://img.shields.io/coveralls/github/LogicalOverflow/music-sync/master.svg?style=flat-square)](https://coveralls.io/github/LogicalOverflow/music-sync?branch=master)
[![Codacy Grade](https://img.shields.io/codacy/grade/d33215d6a3ca41d597ec28c3f06fbf88/master.svg?style=flat-square)](https://www.codacy.com/app/LogicalOverflow/music-sync)
[![GitHub Release Version](https://img.shields.io/github/release/LogicalOverflow/music-sync.svg?style=flat-square)](https://github.com/LogicalOverflow/music-sync/releases/latest)
[![GoDoc Reference](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/LogicalOverflow/music-sync)
[![License](https://img.shields.io/github/license/LogicalOverflow/music-sync.svg?style=flat-square)](https://github.com/LogicalOverflow/music-sync/blob/master/LICENSE)

# Music Sync
A go application to play the same music from multiple devices at once. It works best when all playing devices are similar, to avoid differences in the time it takes the audio to be played. I usually test the timing with two Windows 7 machines, one amd64, one 386 to ensure the timing difference between the devices is small enough that hearing 2 different devices playing music sounds like one.

## Installation
To install do `go get github.com/LogicalOverflow/music-sync/...` or download the executable from the [latest release](https://github.com/LogicalOverflow/music-sync/releases/latest).

## Getting Started
After installing, create a directory named `audio` and put your audio files into it. Also, create a `users.json` containing at least one username and password of your choice:

```json
{
    "username": "password",
    "another username": "another password"
 }
```
Then you can start a local music-sync-server using
`music-sync-server`. By default, this server listens on `127.0.0.1:13333` (`--address`, `--port`) for clients and provides a ssh terminal on `127.0.0.1:13334` (`--ssh-address`, `--ssh-port`) to control the server. By default, the server checks in it's working directory for a file called `id_rsa` to use as a host key (`--host-key-file`). If this file is not found a new host key is generated on every startup. For more options check `music-sync-server --help`.

To start a player use `music-sync-player`. By default this tries to connect to a server at `127.0.0.1:1333` (`--address`, `--port`). For more options check `music-sync-player --help`.

The ssh terminal on the server is used to control the server. The usernames and passwords are read from `users.json` (`--users-file`). You can manage the current playlist, pause and resume playback and set the playback volume for all clients. Theses commands are avilable:
 * `queue filename [position]` - Adds filename to the playlist at position or the end.
 * `remove position` - Removes the song at position from the playlist
 * `jump position` - Jumps to position in the playlist, interrupting the current song
 * `playlist` - Prints the current playlist
 * `pause` - Pauses playback
 * `resume` - Resumes playback
 * `volume volume` - Sets the playback volume for all clients (volume should be between 0 and 1)
 * `help [command]` - Prints all commands or information and usage of command
 * `ls [sub-directory]` - Lists all files in the music (sub-)directory
 * `clear` - Clears the terminal
 * `exit` - Closes the connection
