package main

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"os/signal"
)

// TODO ML Fix strange window behaivor
// TODO ML Should we use SDL sound functions?
func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	fmt.Println("Starting recording, press CTRL+C to abort.")

	// Initialize audio library.
	portaudio.Initialize()
	defer portaudio.Terminate()

	// Initialize visualization library.
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()
	window, err := sdl.CreateWindow("go guitar", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
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
	in := make([]int32, 8192)
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	check(err)
	defer stream.Close()

	// Read sampleRate/len(in) samples per second.
	sampleN := 0
	check(stream.Start())
	for {
		check(stream.Read())

		// Render values in window
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		renderer.SetDrawColor(255, 0, 0, 255)
		f := 800 / float32(len(in))

		amplit := 500000000
		g := 600 / float32(amplit)
		for i := range in {
			x := int(float32(i) * f)
			y := int32(float32(in[x])*g + 300)
			renderer.DrawPoint(int32(x), y)
		}
		renderer.Present()

		// Check if we should exit?
		select {
		case <-sig:
			return
		default:
		}

		sampleN += 1
	}
	check(stream.Stop())
}

// Helper function to check on any error.
func check(err error) {
	if err != nil {
		panic(err)
	}
}
