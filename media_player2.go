package main

import (
	"os"

	"github.com/godbus/dbus/v5/prop"
)

type MediaPlayer2 struct {
	mp    *MocP
	Props map[string]*prop.Prop
}

func NewMediaPlayer2(mp *MocP) *MediaPlayer2 {
	mp2 := &MediaPlayer2{}
	mp2.mp = mp
	mp2.Props = map[string]*prop.Prop{
		"CanQuit":          newProp(true, nil),
		"Fullscreen":       newProp(false, nil),
		"CanSetFullscreen": newProp(false, nil),
		"CanRaise":         newProp(false, nil),
		"HasTrackList":     newProp(false, nil),
		"Identity":         newProp("Media On Console", nil),
	}

	return mp2
}

func (m *MediaPlayer2) Raise() {
}

func (m *MediaPlayer2) Quit() {
	m.mp.Exit()
	os.Exit(0)
}
