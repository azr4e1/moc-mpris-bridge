package main

import (
	"log"
	"os"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
)

type MediaPlayer2 struct {
	mp         *MocP
	conn       *dbus.Conn
	properties *prop.Properties
}

func NewMediaPlayer2(conn *dbus.Conn, mp *MocP) (*MediaPlayer2, error) {
	mp2 := &MediaPlayer2{}
	mp2.mp = mp
	mp2.conn = conn
	props := map[string]*prop.Prop{
		"CanQuit":          newProp(true, nil),
		"Fullscreen":       newProp(false, nil),
		"CanSetFullscreen": newProp(false, nil),
		"CanRaise":         newProp(false, nil),
		"HasTrackList":     newProp(false, nil),
		"Identity":         newProp("Media On Console", nil),
	}

	var err error
	mp2.properties, err = prop.Export(mp2.conn, "/org/mpris/MediaPlayer2", map[string]map[string]*prop.Prop{
		"org.mpris.MediaPlayer2": props,
	})
	if err != nil {
		return nil, err
	}
	log.Println("MediaPlayer2 properties exported")

	return mp2, nil
}

func (m *MediaPlayer2) Raise() {
}

func (m *MediaPlayer2) Quit() {
	m.mp.Exit()
	m.conn.Close()
	os.Exit(0)
}
