package main

import (
	"errors"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
)

type MediaPlayer2Player struct {
	lock       sync.Mutex
	mp         *MocP
	conn       *dbus.Conn
	propValues map[string]any
	properties *prop.Properties
}

func NewMediaPlayer2Player(conn *dbus.Conn, mp *MocP) (*MediaPlayer2Player, error) {
	mp2p := &MediaPlayer2Player{}
	mp2p.lock = sync.Mutex{}
	mp2p.mp = mp
	mp2p.conn = conn

	var err error
	mp2p.propValues, mp2p.properties, err = mp2p.exportProps()
	if err != nil {
		return nil, err
	}
	log.Println("MediaPlayer2.Player properties exported")

	return mp2p, nil
}

func (mp2p *MediaPlayer2Player) exportProps() (map[string]any, *prop.Properties, error) {
	propValues := make(map[string]any)
	propertiesMap := make(map[string]*prop.Prop)

	setProp := func(name string, value any, setter func(*prop.Change) *dbus.Error) {
		propValues[name] = value
		propertiesMap[name] = newProp(value, setter)
	}

	setProp("PlaybackStatus", mp2p.getPlaybackStatus(), nil)
	setProp("LoopStatus", mp2p.getLoopStatus(), mp2p.setLoopStatus)
	setProp("Rate", 1.00, mp2p.setRate)
	setProp("Shuffle", mp2p.getShuffle(), mp2p.setShuffle)
	setProp("Metadata", mp2p.getMetadata(), nil)
	setProp("Volume", mp2p.getVolume(), mp2p.setVolume)
	setProp("MinimumRate", 1.00, nil)
	setProp("MaximumRate", 1.00, nil)
	setProp("CanGoNext", mp2p.getCanGoNext(), nil)
	setProp("CanGoPrevious", mp2p.getCanGoPrevious(), nil)
	setProp("CanPlay", mp2p.getCanPlay(), nil)
	setProp("CanPause", mp2p.getCanPause(), nil)
	setProp("CanSeek", mp2p.getCanSeek(), nil)
	setProp("CanControl", true, nil)

	// The org.freedesktop.DBus.Properties.PropertiesChanged signal is not emitted when this property (Position) changes. We need to set EmitFalse in the prop creation
	posName := "Position"
	posValue := mp2p.getPosition()
	propValues[posName] = posValue
	propertiesMap[posName] = &prop.Prop{
		Value:    posValue,
		Writable: true,
		Emit:     prop.EmitFalse,
		Callback: nil,
	}

	properties, err := prop.Export(mp2p.conn, "/org/mpris/MediaPlayer2", map[string]map[string]*prop.Prop{
		"org.mpris.MediaPlayer2.Player": propertiesMap,
	})

	return propValues, properties, err
}

func (mp2p *MediaPlayer2Player) update() *dbus.Error {
	mp2p.lock.Lock()
	defer mp2p.lock.Unlock()
	err := mp2p.mp.UpdateInfo()
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	for key, value := range mp2p.propValues {
		newVal := mp2p.getCurrVal(key)
		if !reflect.DeepEqual(newVal, value) {
			mp2p.propValues[key] = newVal
			mp2p.properties.Set("org.mpris.MediaPlayer2.Player", key, dbus.MakeVariant(newVal))
			log.Printf("MediaPlayer2.Player.%s was updated\n", key)
		}
	}
	return nil
}

func (mp2p *MediaPlayer2Player) getInfo(key string) any {
	mp2p.lock.Lock()
	val := mp2p.mp.GetInfo(key)
	mp2p.lock.Unlock()

	return val
}

func (mp2p *MediaPlayer2Player) getCurrVal(key string) any {
	switch key {
	case "PlaybackStatus":
		return mp2p.getPlaybackStatus()
	case "LoopStatus":
		return mp2p.getLoopStatus()
	case "Rate":
		return 1.00
	case "Shuffle":
		return mp2p.getShuffle()
	case "Metadata":
		return mp2p.getMetadata()
	case "Volume":
		return mp2p.getVolume()
	case "Position":
		return mp2p.getPosition()
	case "MinimumRate":
		return 1.00
	case "MaximumRate":
		return 1.00
	case "CanGoNext":
		return mp2p.getCanGoNext()
	case "CanGoPrevious":
		return mp2p.getCanGoPrevious()
	case "CanPlay":
		return mp2p.getCanPlay()
	case "CanPause":
		return mp2p.getCanPause()
	case "CanSeek":
		return mp2p.getCanSeek()
	case "CanControl":
		return true
	default:
		return nil
	}
}

// Methods
func (mp2p *MediaPlayer2Player) Next() *dbus.Error {
	if !mp2p.getCanGoNext() {
		log.Println("MediaPlayer2.Player.Next is not allowed")
		return nil
	}
	log.Println("MediaPlayer2.Player.Next was called")
	err := mp2p.mp.Next()
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	time.Sleep(time.Second / 2)
	return mp2p.update()
}

func (mp2p *MediaPlayer2Player) Previous() *dbus.Error {
	if !mp2p.getCanGoPrevious() {
		log.Println("MediaPlayer2.Player.Previous is not allowed")
		return nil
	}
	log.Println("MediaPlayer2.Player.Previous was called")
	err := mp2p.mp.Previous()
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	time.Sleep(time.Second / 2)
	return mp2p.update()
}

func (mp2p *MediaPlayer2Player) Pause() *dbus.Error {
	if !mp2p.getCanPause() {
		log.Println("MediaPlayer2.Player.Pause is not allowed")
		return nil
	}
	log.Println("MediaPlayer2.Player.Pause was called")
	err := mp2p.mp.Pause()
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return mp2p.update()
}

func (mp2p *MediaPlayer2Player) PlayPause() *dbus.Error {
	if !mp2p.getCanPlay() || !mp2p.getCanPause() {
		log.Println("MediaPlayer2.Player.PlayPause is not allowed")
		return nil
	}
	log.Println("MediaPlayer2.Player.PlayPause was called")
	err := mp2p.mp.TogglePause()
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return mp2p.update()
}

func (mp2p *MediaPlayer2Player) Stop() *dbus.Error {
	log.Println("MediaPlayer2.Player.Stop was called")
	err := mp2p.mp.Stop()
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return mp2p.update()
}

func (mp2p *MediaPlayer2Player) Play() *dbus.Error {
	if !mp2p.getCanPlay() {
		log.Println("MediaPlayer2.Player.Play is not allowed")
		return nil
	}
	log.Println("MediaPlayer2.Player.Play was called")
	err := mp2p.mp.Unpause()
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return mp2p.update()
}

func (mp2p *MediaPlayer2Player) Seek(microseconds int64) *dbus.Error {
	if !mp2p.getCanSeek() {
		log.Println("MediaPlayer2.Player.Seek is not allowed")
		return nil
	}
	log.Println("MediaPlayer2.Player.Seek was called")
	seconds := int(microseconds) / 1000000
	err := mp2p.mp.Seek(seconds)
	if err != nil {
		return dbus.MakeFailedError(err)
	}

	current_seconds, ok := mp2p.getInfo(CurrentSec).(int)
	if !ok {
		return nil
	}
	current_microseconds := int64(current_seconds+seconds) * 1000000
	err = mp2p.Seeked(current_microseconds)
	if err != nil {
		return dbus.MakeFailedError(err)
	}

	errUpdate := mp2p.update()
	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

func (mp2p *MediaPlayer2Player) SetPosition(trackId dbus.ObjectPath, microseconds int64) *dbus.Error {
	if !mp2p.getCanSeek() {
		log.Println("MediaPlayer2.Player.SetPosition is not allowed")
		return nil
	}
	log.Println("MediaPlayer2.Player.SetPosition was called")
	seconds := int(microseconds) / 1000000
	err := mp2p.mp.Jump(seconds)
	if err != nil {
		return dbus.MakeFailedError(err)
	}

	err = mp2p.Seeked(microseconds)
	if err != nil {
		return dbus.MakeFailedError(err)
	}

	errUpdate := mp2p.update()
	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

// Signal

func (mp2p *MediaPlayer2Player) Seeked(position int64) error {
	log.Println("MediaPlayer2.Player.Seeked was signalled")
	err := mp2p.conn.Emit("/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2.Player.Seeked", position)

	return err
}

// Properties

func (mp2p *MediaPlayer2Player) getPlaybackStatus() string {
	return mp2p.mp.GetPlaybackStatus()
}

func (mp2p *MediaPlayer2Player) getLoopStatus() string {
	// TODO: implement loop status fetch
	return mp2p.mp.GetLoopStatus()
}

func (mp2p *MediaPlayer2Player) setLoopStatus(change *prop.Change) *dbus.Error {
	value, ok := change.Value.(string)
	if !ok {
		return dbus.MakeFailedError(errors.New("wrong LoopStatus change"))
	}
	switch value {
	case "None":
		err := mp2p.mp.SetRepeat(false)
		if err != nil {
			return dbus.MakeFailedError(err)
		}
	case "Track", "Playlist":
		err := mp2p.mp.SetRepeat(true)
		if err != nil {
			return dbus.MakeFailedError(err)
		}
	}

	return nil
}

func (mp2p *MediaPlayer2Player) setRate(change *prop.Change) *dbus.Error {
	value, ok := change.Value.(float64)
	if !ok {
		return dbus.MakeFailedError(errors.New("wrong Rate change"))
	}
	err := mp2p.mp.SetRate(value)
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

func (mp2p *MediaPlayer2Player) getShuffle() bool {
	// TODO: implement loop status fetch
	return mp2p.mp.GetShuffle()
}

func (mp2p *MediaPlayer2Player) setShuffle(change *prop.Change) *dbus.Error {
	value, ok := change.Value.(bool)
	if !ok {
		return dbus.MakeFailedError(errors.New("wrong Shuffle change"))
	}
	err := mp2p.mp.SetShuffle(value)
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

func (mp2p *MediaPlayer2Player) getMetadata() map[string]any {
	return mp2p.mp.GetMetadata()
}

func (mp2p *MediaPlayer2Player) getVolume() float64 {
	return float64(mp2p.mp.GetVolume()) / 100
}

func (mp2p *MediaPlayer2Player) setVolume(change *prop.Change) *dbus.Error {
	log.Println("MediaPlayer2.Player setVolume was called")
	volume, ok := change.Value.(float64)
	if !ok {
		return dbus.MakeFailedError(errors.New("wrong Volume change"))
	}
	newVol := int(100 * volume)
	err := mp2p.mp.Volume(newVol)

	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

func (mp2p *MediaPlayer2Player) getPosition() int64 {
	return int64(mp2p.mp.GetPosition()) * 1000000
}

func (mp2p *MediaPlayer2Player) getCanGoNext() bool {
	return mp2p.mp.CanGoNext()
}

func (mp2p *MediaPlayer2Player) getCanGoPrevious() bool {
	return mp2p.mp.CanGoPrev()
}

func (mp2p *MediaPlayer2Player) getCanPlay() bool {
	return mp2p.mp.CanPlay()
}

func (mp2p *MediaPlayer2Player) getCanPause() bool {
	return mp2p.mp.CanPause()
}

func (mp2p *MediaPlayer2Player) getCanSeek() bool {
	return mp2p.mp.CanSeek()
}
