package main

import (
	"flag"
	"fmt"
	"github.com/hajimehoshi/oto/v2"
	"io/ioutil"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	audioPathPrefix = "audio/"
)

var pitchMap = map[string]string{
	"1--": audioPathPrefix + "ll1.mp3",
	"2--": audioPathPrefix + "ll2.mp3",
	"3--": audioPathPrefix + "ll3.mp3",
	"4--": audioPathPrefix + "ll4.mp3",
	"5--": audioPathPrefix + "ll5.mp3",
	"6--": audioPathPrefix + "ll6.mp3",
	"7--": audioPathPrefix + "ll7.mp3",
	"1-": audioPathPrefix + "l1.mp3",
	"2-": audioPathPrefix + "l2.mp3",
	"3-": audioPathPrefix + "l3.mp3",
	"4-": audioPathPrefix + "l4.mp3",
	"5-": audioPathPrefix + "l5.mp3",
	"6-": audioPathPrefix + "l6.mp3",
	"7-": audioPathPrefix + "l7.mp3",
	"1": audioPathPrefix + "m1.mp3",
	"2": audioPathPrefix + "m2.mp3",
	"3": audioPathPrefix + "m3.mp3",
	"4": audioPathPrefix + "m4.mp3",
	"5": audioPathPrefix + "m5.mp3",
	"6": audioPathPrefix + "m6.mp3",
	"7": audioPathPrefix + "m7.mp3",
	"1+": audioPathPrefix + "h1.mp3",
	"2+": audioPathPrefix + "h2.mp3",
	"3+": audioPathPrefix + "h3.mp3",
	"4+": audioPathPrefix + "h4.mp3",
	"5+": audioPathPrefix + "h5.mp3",
	"6+": audioPathPrefix + "h6.mp3",
	"7+": audioPathPrefix + "h7.mp3",
	"1++": audioPathPrefix + "hh1.mp3",
	"2++": audioPathPrefix + "hh2.mp3",
	"3++": audioPathPrefix + "hh3.mp3",
	"4++": audioPathPrefix + "hh4.mp3",
	"5++": audioPathPrefix + "hh5.mp3",
	"6++": audioPathPrefix + "hh6.mp3",
	"7++": audioPathPrefix + "hh7.mp3",
}

var animationMap = map[string]string{
	"0": "_",
	"1": "▁",
	"2": "▂",
	"3": "▃",
	"4": "▄",
	"5": "▅",
	"6": "▆",
	"7": "▇",
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "array var"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var pianoNotes arrayFlags
var interval time.Duration

func loadText(path string) (texts []string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic("read file error: " + err.Error())
	}
	texts = strings.Split(string(data), " ")
	return
}

func main() {
	flag.Var(&pianoNotes, "notes", "load piano notes file")
	flag.DurationVar(&interval, "interval", 180*time.Millisecond, "audio pitch interval time")
	flag.Parse()

	fmt.Println("load file: ", pianoNotes)
	fmt.Println("interval time: ", interval)
	if len(pianoNotes) == 0 {
		panic("piano notes load failed")
	}

	c, ready, err := oto.NewContext(48000, 2, 2)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-ready

	wg := &sync.WaitGroup{}
	wg.Add(len(pianoNotes))
	for _, path := range pianoNotes {
		go loadPiano(c, path, wg)
	}
	go loadAnimation(pianoNotes[0])
	//go loadPiano(c, "./resource/notes/起风了_180.notes")
	//go loadPiano(c, "./resource/notes/起风了_180.accompaniments")

	wg.Wait()

}

func loadPiano(c *oto.Context, path string, wg *sync.WaitGroup) {
	defer wg.Done()
	times := interval
	f, preAudio := decoderMp3("audio/test.mp3")
	go play(c, preAudio, f)
	time.Sleep(1 * time.Second)

	texts := loadText(path)
	for _, text := range texts {
		text = strings.Replace(text, "\n", "", -1)
		text = strings.Replace(text, "\r", "", -1)
		text = strings.Replace(text, "\t", "", -1)
		text = strings.TrimSpace(text)
		if utf8.RuneCountInString(text) < 1 {
			continue
		}
		if value, ok := pitchMap[text]; ok {
			f, d := decoderMp3(value)
			go play(c, d, f)
			time.Sleep(times/2)
		} else if text == "0" {
			time.Sleep(times/2)
		} else {
			continue
		}
		time.Sleep(times/2)
		times = interval
	}
}

func loadAnimation(path string) {
	time.Sleep(1*time.Second)
	texts := loadText(path)
	for _, text := range texts {
		if utf8.RuneCountInString(text) < 1 {
			continue
		}
		text = strings.Replace(text, "-", "", -1)
		text = strings.Replace(text, "+", "", -1)

		if value, ok := animationMap[text]; ok {
			fmt.Print(value)
			fmt.Print(" ")
		} else {
			fmt.Println()
		}
		time.Sleep(interval)
	}
}