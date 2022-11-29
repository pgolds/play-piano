package main

import (
	"embed"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"io/fs"
	"time"
)

//go:embed resource/audio/**
var resource embed.FS

func decoderMp3(path string) (fs.File, *mp3.Decoder) {
	sub, err := fs.Sub(resource, "resource")
	if err != nil {
		panic(err)
	}
	f, err := sub.Open(path)
	if err != nil {
		panic("reading my-file.mp3 failed: " + err.Error())
	}

	d, err := mp3.NewDecoder(f)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}
	return f, d
}

func play(c *oto.Context, d *mp3.Decoder, f fs.File) {

	// Create a new 'player' that will handle our sound. Paused by default.
	player := c.NewPlayer(d)


	player.(oto.BufferSizeSetter).SetBufferSize(4096)

	// Play starts playing the sound and returns without waiting for it (Play() is async).
	player.Play()

	// We can wait for the sound to finish playing using something like this
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	err := player.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())
	}
	f.Close()
}
