package drum

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

var (
	hihatSound  = "hihat"
	kickSound   = "kick"
	snareSound  = "snare"
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

func init() {
	log.Println("Drum initialized...")
	sampleRate := beep.SampleRate(44100)
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))
}

type Drum struct {
	TimeDrum int
}

func (d Drum) Hihat(rhythm string) {
	d.playbeat("prato", rhythm, hihatSound, colorYellow)
}

func (d Drum) Kick(rhythm string) {
	d.playbeat("bumbo", rhythm, kickSound, colorRed)
}

func (d Drum) Snare(rhythm string) {
	d.playbeat("caixa", rhythm, snareSound, colorGreen)
}

func (d Drum) playbeat(name, beats, sound, color string) {
	ticker := time.NewTicker(time.Duration(d.TimeDrum) * time.Millisecond)
	defer ticker.Stop()

	runes := []rune(beats)
	count := 0

	for count < len(runes) {
		select {
		case _ = <-ticker.C:
			if runes[count] == rune('x') {
				playSound(loadSound(sound))
				fmt.Println(string(color), "P"+strconv.Itoa(count+1)+" -> "+name, string(colorReset))
			}
		}
		count++
	}
}

func loadSound(name string) beep.StreamSeekCloser {
	a, err := filepath.Abs(`./assets/audio/` + name + ".wav")

	if err != nil {
		log.Println(err)
	}

	f, err := os.Open(a)

	if err != nil {
		log.Println(err)
	}

	stream, _, err := wav.Decode(f)

	if err != nil {
		log.Println(err)
	}

	return stream
}

func playSound(sound beep.StreamSeeker) {
	done := make(chan bool)

	speaker.Play(beep.Seq(sound, beep.Callback(func() {
		done <- true
	})))

	<-done

	sound.(beep.StreamSeekCloser).Close()
}
