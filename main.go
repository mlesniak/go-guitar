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

	// Initialize visualization library.
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Go guitar", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return
	}
	defer renderer.Destroy()

	// Create buffer for sampling data and open the stream.
	in := make([]int32, bufferSize)
	stream, err := portaudio.OpenDefaultStream(1, 0, sampleRate, len(in), in)
	check(err)
	defer stream.Close()

	check(stream.Start())
	f := width / float32(len(in))
	amplit := 500000000
	g := height / float32(amplit)
	for {
		check(stream.Read())

		// Render values in window
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		renderer.SetDrawColor(255, 0, 0, 255)
		for i := range in {
			x := int(float32(i) * f)
			y := int32(float32(in[x])*g + height/2)
			renderer.DrawPoint(int32(x), y)
		}
		renderer.Present()

		// See https://stackoverflow.com/questions/39637824/border-titlebar-not-properly-displaying-in-sdl-osx
		var event sdl.Event
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.KeyboardEvent:
				if t.Keysym.Sym == 27 {
					// Quit on Escape key.
					return
				}
			}
		}

		// Check if we should exit?
		select {
		case <-sig:
			return
		default:
		}
	}
	check(stream.Stop())
}

// Helper function to check on any error.
func check(err error) {
	if err != nil {
		panic(err)
	}
}
