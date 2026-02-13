package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/godbus/dbus/v5"
)

func MPRISLoop(name string) error {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Println("DBus connection created")

	mp, err := NewMocP()
	if err != nil {
		return err
	}
	log.Println("MocP instance initialized")

	mp2, err := NewMediaPlayer2(conn, mp)
	if err != nil {
		return err
	}
	log.Println("MediaPlayer2 instance created")

	mp2p, err := NewMediaPlayer2Player(conn, mp)
	if err != nil {
		return err
	}
	log.Println("MediaPlayer2.Player instance created")

	// Register name
	reply, err := conn.RequestName(fmt.Sprintf("org.mpris.MediaPlayer2.%s", name), dbus.NameFlagReplaceExisting)
	if err != nil {
		return err
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		return errors.New("Name is already taken")
	}
	log.Printf("%s name successfully registered\n", name)

	err = conn.Export(mp2, "/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2")
	if err != nil {
		return err
	}
	log.Println("MediaPlayer2 interface exported")

	err = conn.Export(mp2p, "/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2.Player")
	if err != nil {
		return err
	}
	log.Println("MediaPlayer2.Player interface exported")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)

	log.Println("Starting loop...")

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case dbusMethod := <-mp2p.commands:
			// Dbus methods asked to do something
			dbusMethod.result <- dbusMethod.action()
			err := mp2p.update()
			if err != nil {
				return err
			}
		case <-ticker.C:
			// poll every second
			err := mp2p.update()
			if err != nil {
				return err
			}
		case <-c:
			log.Println("Interruption...")
			return nil
		}
	}
}
