package main

import (
	"errors"
	"math"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
)

type MocP struct {
	metadata map[string]any
}

const (
	State       = "State"
	File        = "File"
	Title       = "Title"
	Artist      = "Artist"
	SongTitle   = "SongTitle"
	Album       = "Album"
	TotalTime   = "TotalTime"
	TimeLeft    = "TimeLeft"
	TotalSec    = "TotalSec"
	CurrentTime = "CurrentTime"
	CurrentSec  = "CurrentSec"
	Bitrate     = "Bitrate"
	AvgBitrate  = "AvgBitrate"
	Rate        = "Rate"
)

var mocpInfoKeys = map[string]bool{
	State:       true,
	File:        true,
	Title:       true,
	Artist:      true,
	SongTitle:   true,
	Album:       true,
	TotalTime:   true,
	TimeLeft:    true,
	TotalSec:    true,
	CurrentTime: true,
	CurrentSec:  true,
	Bitrate:     true,
	AvgBitrate:  true,
	Rate:        true,
}

func NewMocP() (*MocP, error) {
	metadata := make(map[string]any)
	mp := &MocP{metadata: metadata}
	err := mp.UpdateInfo()
	if err != nil {
		return nil, err
	}

	return mp, nil
}

func (mp *MocP) Append(files []string) error {
	if mp == nil {
		return nil
	}
	args := []string{"-a"}
	args = append(args, files...)
	cmd := exec.Command("mocp", args...)

	return cmd.Run()
}

func (mp *MocP) Enqueue(files []string) error {
	if mp == nil {
		return nil
	}
	args := []string{"-q"}
	args = append(args, files...)
	cmd := exec.Command("mocp", args...)

	return cmd.Run()
}

func (mp *MocP) ToggleShuffle() error {
	if mp == nil {
		return nil
	}
	cmd := exec.Command("mocp", "-t", "shuffle")
	return cmd.Run()
}

func (mp *MocP) ToggleAutoNext() error {
	if mp == nil {
		return nil
	}
	cmd := exec.Command("mocp", "-t", "autonext")
	return cmd.Run()
}

func (mp *MocP) ToggleRepeat() error {
	if mp == nil {
		return nil
	}
	cmd := exec.Command("mocp", "-t", "repeat")
	return cmd.Run()
}

func (mp *MocP) SetShuffle(on bool) error {
	if mp == nil {
		return nil
	}
	var cmd *exec.Cmd
	if on {
		cmd = exec.Command("mocp", "-o", "shuffle")
	} else {
		cmd = exec.Command("mocp", "-u", "shuffle")
	}
	return cmd.Run()
}

func (mp *MocP) SetAutoNext(on bool) error {
	if mp == nil {
		return nil
	}
	var cmd *exec.Cmd
	if on {
		cmd = exec.Command("mocp", "-o", "autonext")
	} else {
		cmd = exec.Command("mocp", "-u", "autonext")
	}
	return cmd.Run()
}

func (mp *MocP) SetRepeat(on bool) error {
	if mp == nil {
		return nil
	}
	var cmd *exec.Cmd
	if on {
		cmd = exec.Command("mocp", "-o", "repeat")
	} else {
		cmd = exec.Command("mocp", "-u", "repeat")
	}
	return cmd.Run()
}

func (mp *MocP) Clear() error {
	if mp == nil {
		return nil
	}
	cmd := exec.Command("mocp", "-c")
	return cmd.Run()
}

func (mp *MocP) Previous() error {
	if mp == nil {
		return nil
	}
	cmd := exec.Command("mocp", "-r")
	return cmd.Run()
}

func (mp *MocP) Next() error {
	if mp == nil {
		return nil
	}
	cmd := exec.Command("mocp", "-f")
	return cmd.Run()
}

func (mp *MocP) Stop() error {
	if mp == nil {
		return nil
	}
	cmd := exec.Command("mocp", "-s")
	return cmd.Run()
}

func (mp *MocP) Exit() error {
	if mp == nil {
		return nil
	}
	cmd := exec.Command("mocp", "-x")
	return cmd.Run()
}

func (mp *MocP) Unpause() error {
	if mp == nil {
		return nil
	}
	cmd := exec.Command("mocp", "-U")
	return cmd.Run()
}

func (mp *MocP) Play() error {
	if mp == nil {
		return nil
	}
	cmd := exec.Command("mocp", "-p")
	return cmd.Run()
}

func (mp *MocP) TogglePause() error {
	if mp == nil {
		return nil
	}
	state := mp.GetPlaybackStatus()
	switch state {
	case "Playing":
		err := mp.Pause()
		if err != nil {
			return err
		}
	case "Paused":
		err := mp.Unpause()
		if err != nil {
			return err
		}
	}
	return nil
}

func (mp *MocP) GetLoopStatus() string {
	return "None"
}

func (mp *MocP) GetShuffle() bool {
	return false
}

func (mp *MocP) Seek(seconds int) error {
	if mp == nil {
		return nil
	}
	cmd := exec.Command("mocp", "--seek", strconv.Itoa(seconds))
	return cmd.Run()
}

func (mp *MocP) Volume(val int) error {
	if mp == nil {
		return nil
	}
	var volume int
	switch {
	case val < 0:
		volume = 0
	case val > 100:
		volume = 100
	default:
		volume = val
	}

	cmd := exec.Command("mocp", "--volume", strconv.Itoa(volume))
	return cmd.Run()
}

func (mp *MocP) Jump(seconds int) error {
	if mp == nil {
		return nil
	}
	totSec, ok := mp.metadata[TotalSec]
	if !ok {
		return nil
	}
	if seconds < 0 || seconds > totSec.(int) {
		return nil
	}
	cmd := exec.Command("mocp", "--jump", strconv.Itoa(seconds)+"s")
	return cmd.Run()
}

func (mp *MocP) Pause() error {
	if mp == nil {
		return nil
	}
	cmd := exec.Command("mocp", "-P")
	return cmd.Run()
}

func (mp *MocP) GetPlaybackStatus() string {
	if mp == nil {
		return ""
	}
	state, ok := mp.metadata[State]
	if !ok {
		return "Stopped"
	}
	switch state {
	case "STOP":
		return "Stopped"
	case "PLAY":
		return "Playing"
	case "PAUSE":
		return "Paused"
	default:
		return "Stopped"
	}
}

// func (mp *MocP) GetLoopStatus(val string)
// func (mp *MocP) SetLoopStatus(val string)

func (mp *MocP) SetRate(val float64) error {
	if mp == nil {
		return nil
	}
	if val == 0 {
		return mp.Pause()
	}

	return nil
}

// func (mp *MocP) GetShuffle
// func (mp *MocP) SetShuffle

func (mp *MocP) GetMetadata() map[string]any {
	if mp == nil {
		return nil
	}
	metadata := make(map[string]any)

	if val, ok := mp.metadata[TotalSec]; ok {
		metadata["mpris:length"] = int64(val.(int)) * 1000000
	}
	if val, ok := mp.metadata[File]; ok {
		metadata["xesam:url"] = val
		metadata["mpris:artUrl"] = val
	}
	if val, ok := mp.metadata[SongTitle]; ok {
		metadata["xesam:title"] = val
	}
	if val, ok := mp.metadata[Artist]; ok {
		metadata["xesam:artist"] = val
	}
	if val, ok := mp.metadata[Album]; ok {
		metadata["xesam:album"] = val
	}
	metadata["mpris:trackid"] = dbus.ObjectPath("/org/moc_mpris_bridge/track/1")
	// TODO: implement musicbrainz album art fetch
	// if val, ok := mp.metadata[]; ok {
	// 	metadata["mpris:artUrl"] = val
	// }

	return metadata
}

func (mp *MocP) GetVolume() int {
	if mp == nil {
		return 0
	}
	vol, err := amixerGetVolume()
	if err != nil {
		return 0
	}
	return int(vol)
}

func (mp *MocP) GetPosition() int {
	if mp == nil {
		return 0
	}
	currSec, ok := mp.metadata[CurrentSec].(int)
	if !ok {
		return 0
	}
	return currSec
}

func (mp *MocP) CanGoNext() bool {
	if mp == nil {
		return false
	}
	if _, ok := mp.metadata[SongTitle]; ok {
		return true
	}
	return false
}

func (mp *MocP) CanGoPrev() bool {
	if mp == nil {
		return false
	}
	if _, ok := mp.metadata[SongTitle]; ok {
		return true
	}
	return false
}

func (mp *MocP) CanPlay() bool {
	if mp == nil {
		return false
	}
	if _, ok := mp.metadata[SongTitle]; ok {
		return true
	}
	return false
}

func (mp *MocP) CanPause() bool {
	if mp == nil {
		return false
	}
	if _, ok := mp.metadata[SongTitle]; ok {
		return true
	}
	return false
}

func (mp *MocP) CanSeek() bool {
	if mp == nil {
		return false
	}
	if _, ok := mp.metadata[CurrentSec]; ok {
		return true
	}
	return false
}

func (mp *MocP) GetInfo(key string) any {
	val, ok := mp.metadata[key]
	if !ok {
		return nil
	}

	return val
}

func (mp *MocP) UpdateInfo() error {
	if mp == nil {
		return errors.New("must initialize mocp")
	}
	cmd := exec.Command("mocp", "-i")
	data, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	// clean metadata
	clear(mp.metadata)
	lines := strings.SplitSeq(string(data), "\n")

	for l := range lines {
		pairs := strings.SplitN(l, ":", 2)
		if len(pairs) != 2 {
			continue
		}
		key, val := pairs[0], pairs[1]
		val = strings.TrimSpace(val)
		if _, ok := mocpInfoKeys[key]; ok {
			switch key {
			case "TotalTime", "TimeLeft", "CurrentTime":
				durationVal, err := parseDuration(val)
				if err != nil {
					return err
				}
				mp.metadata[key] = durationVal
			case "TotalSec", "CurrentSec":
				secondVal, err := strconv.Atoi(val)
				if err != nil {
					return err
				}
				mp.metadata[key] = secondVal
			default:
				mp.metadata[key] = val
			}
		}
	}
	return nil
}

func parseDuration(duration string) (time.Duration, error) {
	// expects at most 3 elements HH:MM:SS
	parts := strings.Split(duration, ":")
	seconds := 0
	length := len(parts) - 1
	for i := length; i >= 0; i-- {
		val := parts[i]
		secVal, err := strconv.Atoi(val)
		if err != nil {
			return 0, err
		}
		seconds += secVal * int(math.Pow(60.0, float64(length)-float64(i)))
	}
	secondDuration := time.Second * time.Duration(seconds)

	return secondDuration, nil
}

func amixerGetVolume() (float64, error) {
	cmd := exec.Command("amixer", "get", "Master")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}
	reg, err := regexp.Compile(`\[([0-9]*)%\]`)
	if err != nil {
		return 0, err
	}
	matches := reg.FindAllStringSubmatch(string(out), -1)
	if matches == nil {
		return 0, errors.New("couldn't get volume from ALSA")
	}
	volumes := 0
	// average out the speakers volume
	for _, m := range matches {
		submatch := strings.TrimSpace(m[1])
		vol, err := strconv.Atoi(submatch)
		if err != nil {
			return 0, err
		}
		volumes += vol
	}

	return float64(volumes / len(matches)), nil
}
