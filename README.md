[![Build Status](https://img.shields.io/travis/LogicalOverflow/music-sync/master.svg?style=flat-square)](https://travis-ci.org/LogicalOverflow/music-sync)
[![Coverage Status](https://img.shields.io/coveralls/github/LogicalOverflow/music-sync/master.svg?style=flat-square)](https://coveralls.io/github/LogicalOverflow/music-sync?branch=master)
[![Codacy Grade](https://img.shields.io/codacy/grade/8549c99721aa415788943169b621f8de/master.svg?style=flat-square)](https://www.codacy.com/app/LogicalOverflow/music-sync)
[![Code Climate Issues](https://img.shields.io/codeclimate/issues/github/LogicalOverflow/music-sync.svg?style=flat-square)](https://codeclimate.com/github/LogicalOverflow/music-sync/issues)
[![Code Climate Maintainability](https://img.shields.io/codeclimate/maintainability/LogicalOverflow/music-sync.svg?style=flat-square)](https://codeclimate.com/github/LogicalOverflow/music-sync)
[![GitHub Release Version](https://img.shields.io/github/release/LogicalOverflow/music-sync.svg?style=flat-square)](https://github.com/LogicalOverflow/music-sync/releases/latest)
[![GoDoc Reference](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/LogicalOverflow/music-sync)
[![License](https://img.shields.io/github/license/LogicalOverflow/music-sync.svg?style=flat-square)](https://github.com/LogicalOverflow/music-sync/blob/master/LICENSE)

# Music Sync
A go application to play the same music on multiple devices at once. It works best when all playing devices are similar, to avoid differences in the time it takes the audio to be played. I usually test the timing with two Windows 7 machines, one x64, one x86 to ensure the timing difference between the devices is small enough that hearing 2 different devices playing music sounds like one.

## Installation
To install do `go get github.com/LogicalOverflow/music-sync/...` or download the executable from the [latest release](https://github.com/LogicalOverflow/music-sync/releases/latest).

## Getting Started
After installing, create a directory named `audio` and put your audio files into it. Also, create a `users.json` containing at least one username and password/public key of your choice:

```json
{
    "username": {
        "password": "a password",
        "pubKey": "a public key formatted for use in the SSH wire protocol (RFC 4253, section 6.6)"
    }
 }
```
If you want, you can add lyrics information to your songs. To add lyrics to a song called `song.mp3`, create a file called `song.mp3.json` next to the `mp3` file containing the lyrics:

```json
[
    [
        {
            "timestamp": 1234,
            "caption": "The first "
        },
        {
            "timestamp": 5678,
            "caption": "line of lyrics."
        }
    ],
    [
        {
            "timestamp": 9101,
            "caption": "The second "
        },
        {
            "timestamp": 10112,
            "caption": "line."
        }
    ]
]
```
Here, timestamps are in milliseconds from the start of the song, each array describes a line of lyrics and each object in those arrays a word/phrase/syllable in the lyrics.

Then you can start a local music-sync-server using
`music-sync-server`. By default, this server listens on `127.0.0.1:13333` (`--address`, `--port`) for clients and provides a ssh terminal on `127.0.0.1:13334` (`--ssh-address`, `--ssh-port`) to control the server. By default, the server checks in it's working directory for a file called `id_rsa` to use as a host key (`--host-key-file`). If this file is not found a new host key is generated on every startup. For more options check `music-sync-server --help`.

To start a player use `music-sync-player`. By default this tries to connect to a server at `127.0.0.1:1333` (`--address`, `--port`). For more options check `music-sync-player --help`.

To get information about the current song playing and lyrics (if provided) in a terminal UI, you can use `music-sync-infoer`. By default this tries to connect to a server at  `127.0.0.1:1333` (`--address`, `--port`). For more options check `music-sync-infoer --help`.

The ssh terminal on the server is used to control the server. The usernames and passwords are read from `users.json` (`--users-file`). You can manage the current playlist, pause and resume playback and set the playback volume for all clients. These commands are available:
 * `queue filename [position]` - Adds filename to the playlist at position or the end. You can use glob patterns to add multiple files.
 * `remove position` - Removes the song at position from the playlist
 * `jump position` - Jumps to position in the playlist, interrupting the current song
 * `playlist` - Prints the current playlist
 * `pause` - Pauses playback
 * `resume` - Resumes playback
 * `volume volume` - Sets the playback volume for all clients (volume should be between 0 and 1)
 * `help [command]` - Prints all commands or information and usage of command
 * `ls [sub-directory]` - Lists all songs in the music (sub-)directory
 * `clear` - Clears the terminal
 * `exit` - Closes the connection
 
Spaces in commands can be escaped using `\ `. To escape a backslash before a space use `\\ `, otherwise the backslash does not need to be escaped.
