package main

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"os/signal"
)

// Global constants for visualization
const width = 800
const height = 800

// Factors to shrink the measured data with large integers into the window.
const widthFactor = width / float32(bufferSize)
const maxAmplitude = 500000000
const heightFactor = height / float32(maxAmplitude)

// Global constants for audio recording
const bufferSize = 1024
const sampleRate = 88200

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	fmt.Println("Starting recording, press CTRL+C to abort.")

	// Initialize audio library.
	portaudio.Initialize()
	defer portaudio.Terminate()

	// Initialize graphics part.
	check(sdl.Init(sdl.INIT_VIDEO))
	defer sdl.Quit()
	window, err := sdl.CreateWindow("Go guitar", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width, height, sdl.WINDOW_SHOWN)
	check(err)
	defer window.Destroy()
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	defer renderer.Destroy()

	// Create buffer for sampling data and open the stream.
	in := make([]int32, bufferSize)
	stream, err := portaudio.OpenDefaultStream(1, 0, sampleRate, len(in), in)
	defer stream.Close()

	check(stream.Start())
	for {
		check(stream.Read())
		drawVoice(renderer, in)

		if checkForExit(sig) {
			return
		}
	}
	check(stream.Stop())
}

func drawVoice(renderer *sdl.Renderer, in []int32) {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	renderer.SetDrawColor(255, 0, 0, 255)
	for i := range in {
		x := int(float32(i) * widthFactor)
		y := int32(float32(in[x])*heightFactor + height/2)
		renderer.DrawPoint(int32(x), y)
	}
	renderer.Present()
}

func checkForExit(sig chan os.Signal) bool {
	// See https://stackoverflow.com/questions/39637824/border-titlebar-not-properly-displaying-in-sdl-osx
	var event sdl.Event
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return true
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == 27 {
				// Quit on Escape key.
				return true
			}
		}
	}
	// Check if we should exit?
	select {
	case <-sig:
		return true
	default:
	}
	return false
}

// Helper function to check on any error.
func check(err error) {
	if err != nil {
		panic(err)
	}
}
