package main

import (
	"log"

	"github.com/godbus/dbus/v5"
)

func main() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	mp, err := NewMocP()
	if err != nil {
		log.Fatal(err)
	}
	mp2, err := NewMediaPlayer2(conn, mp)
	if err != nil {
		log.Fatal(err)
	}
	mp2p, err := NewMediaPlayer2Player(conn, mp)
	if err != nil {
		log.Fatal(err)
	}

	reply, err := conn.RequestName("org.mpris.MediaPlayer2.mocp-mpris-bridge", dbus.NameFlagReplaceExisting)
	if err != nil {
		log.Fatal(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		log.Fatal("Name already taken")
	}

	err = conn.Export(mp2, "/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2")
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Export(mp2p, "/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2.Player")
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
