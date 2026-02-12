# moc-mpris-bridge

A lightweight D-Bus bridge that implements the [MPRIS](https://specifications.freedesktop.org/mpris-spec/latest/) (MediaPlayer2) interface for [MOC (Music on Console)](https://moc.daper.net/).

This allows MOC to integrate with desktop environments and media controllers that use the standard MPRIS protocol — enabling media key support, playback info in status bars, and control through tools like `playerctl`.

## Features

- Playback control (play, pause, stop, next, previous)
- Seek and position tracking
- Metadata (title, artist, album, duration)
- Volume control via ALSA
- Shuffle and repeat mode support
- Runs as a systemd user service

## Requirements

- [MOC](https://moc.daper.net/) (`mocp`)
- D-Bus
- ALSA (`amixer`) for volume control
- Go 1.25+ (build only)

## Installation

### From AUR

```sh
yay -S moc-mpris-bridge
```

### From source

```sh
go install github.com/azre1/moc-mpris-bridge@latest
```

Or clone and build manually:

```sh
git clone https://github.com/azre1/moc-mpris-bridge.git
cd moc-mpris-bridge
go build -o moc-mpris-bridge .
```

### From releases

Pre-built binaries for linux/amd64 and linux/arm64 are available on the [GitHub Releases](https://github.com/azr4e1/moc-mpris-bridge/releases) page.

## Usage

Start the bridge directly:

```sh
moc-mpris-bridge
```

### Systemd service

Copy the service file and enable it as a user service:

```sh
sudo cp moc-mpris-bridge /usr/bin/
cp moc-mpris-bridge.service ~/.config/systemd/user/
systemctl --user enable --now moc-mpris-bridge
```

The service will automatically restart if it exits, and starts after D-Bus is available.

The service file is available immediately if installing from AUR.

## How it works

The bridge polls MOC's status once per second via `mocp -i` and exposes the state over D-Bus under the name `org.mpris.MediaPlayer2.mocp-mpris-bridge`. It implements both the `org.mpris.MediaPlayer2` and `org.mpris.MediaPlayer2.Player` interfaces, so any MPRIS-aware client can discover and control MOC.


# TODO

- [ ] Improve race conditions when reading mocp info
- [X] get track art through file metadata
