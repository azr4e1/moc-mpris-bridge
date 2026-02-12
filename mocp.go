package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dhowden/tag"
	"github.com/godbus/dbus/v5"
)

type MocP struct {
	lock     sync.Mutex
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
	ArtURI      = "ArtURI"
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
	mp.lock = sync.Mutex{}
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
	// We use the unsage version of getPlaybackStatus
	// to avoid a TOCTOU bug, when checking and
	// acting are not a single atomic action
	mp.lock.Lock()
	defer mp.lock.Unlock()
	state := mp.getPlaybackStatus()
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
	mp.lock.Lock()
	totSec, ok := mp.getInfo(TotalSec)
	mp.lock.Unlock()
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

func (mp *MocP) SafeGetPlaybackStatus() string {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	val := mp.getPlaybackStatus()
	return val
}

func (mp *MocP) getPlaybackStatus() string {
	if mp == nil {
		return ""
	}
	state, ok := mp.getInfo(State)
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
	mp.lock.Lock()
	defer mp.lock.Unlock()
	metadata := make(map[string]any)

	if val, ok := mp.getInfo(TotalSec); ok {
		metadata["mpris:length"] = int64(val.(int)) * 1000000
	}
	if val, ok := mp.getInfo(File); ok {
		metadata["xesam:url"] = val
	}
	if val, ok := mp.getInfo(ArtURI); ok {
		metadata["mpris:artUrl"] = val
	}
	if val, ok := mp.getInfo(SongTitle); ok {
		metadata["xesam:title"] = val
	}
	if val, ok := mp.getInfo(Artist); ok {
		metadata["xesam:artist"] = val
	}
	if val, ok := mp.getInfo(Album); ok {
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
	mp.lock.Lock()
	defer mp.lock.Unlock()
	currSecAny, ok := mp.getInfo(CurrentSec)
	if !ok {
		return 0
	}
	currSec, ok := currSecAny.(int)
	if !ok {
		return 0
	}
	return currSec
}

func (mp *MocP) CanGoNext() bool {
	if mp == nil {
		return false
	}
	mp.lock.Lock()
	defer mp.lock.Unlock()
	if _, ok := mp.getInfo(SongTitle); ok {
		return true
	}
	return false
}

func (mp *MocP) CanGoPrev() bool {
	if mp == nil {
		return false
	}
	mp.lock.Lock()
	defer mp.lock.Unlock()
	if _, ok := mp.getInfo(SongTitle); ok {
		return true
	}
	return false
}

func (mp *MocP) CanPlay() bool {
	if mp == nil {
		return false
	}
	mp.lock.Lock()
	defer mp.lock.Unlock()
	if _, ok := mp.getInfo(SongTitle); ok {
		return true
	}
	return false
}

func (mp *MocP) CanPause() bool {
	if mp == nil {
		return false
	}
	mp.lock.Lock()
	defer mp.lock.Unlock()
	if _, ok := mp.getInfo(SongTitle); ok {
		return true
	}
	return false
}

func (mp *MocP) CanSeek() bool {
	if mp == nil {
		return false
	}
	mp.lock.Lock()
	defer mp.lock.Unlock()
	if _, ok := mp.getInfo(CurrentSec); ok {
		val, ok := mp.getInfo(State)
		if ok && val == "PLAY" {
			return true
		}
	}
	return false
}

func (mp *MocP) getInfo(key string) (any, bool) {
	val, ok := mp.metadata[key]
	return val, ok
}

func (mp *MocP) SafeGetInfo(key string) (any, bool) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	val, ok := mp.metadata[key]
	return val, ok
}

func (mp *MocP) UpdateInfo() error {
	if mp == nil {
		return errors.New("must initialize mocp")
	}
	mp.lock.Lock()
	defer mp.lock.Unlock()
	cmd := exec.Command("mocp", "-i")
	data, err := cmd.CombinedOutput()
	if err != nil {
		// mocp crashed, treat it as stopped
		clear(mp.metadata)
		return nil
	}
	// clean metadata, but keep a copy of file and artURI to avoid reencoding
	file, artURI := mp.cacheArtURI()
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
			case File:
				// if file is the same, don't recalculate artwork
				if val == file {
					mp.metadata[ArtURI] = artURI
				} else {
					mp.metadata[ArtURI] = retrieveArtworkDataURI(val)
				}
				mp.metadata[File] = val
			case TotalTime, TimeLeft, CurrentTime:
				durationVal, err := parseDuration(val)
				if err != nil {
					return err
				}
				mp.metadata[key] = durationVal
			case TotalSec, CurrentSec:
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

func (mp *MocP) cacheArtURI() (string, string) {
	if mp == nil {
		return "", ""
	}
	var file, artURI string
	fileAny, ok1 := mp.getInfo(File)
	artURIAny, ok2 := mp.getInfo(ArtURI)

	if ok1 && ok2 {
		file, ok1 = fileAny.(string)
		artURI, ok2 = artURIAny.(string)
	}

	if ok1 && ok2 {
		return file, artURI
	}

	return "", ""
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

func retrieveArtworkDataURI(file string) string {
	fd, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer fd.Close()

	m, err := tag.ReadFrom(fd)
	if err != nil {
		return ""
	}

	pic := m.Picture()
	if pic == nil {
		return ""
	}

	base64Encoding := base64.StdEncoding.EncodeToString(pic.Data)

	uri := fmt.Sprintf("data:%s;base64,%s", pic.MIMEType, base64Encoding)

	return uri
}
