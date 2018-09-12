package main

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"os/signal"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	fmt.Println("Starting recording, press CTRL+C to abort.")

	// Initialize visualization library.
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// Initialize audio library.
	portaudio.Initialize()
	defer portaudio.Terminate()

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
	in := make([]int32, 1024)
	stream, err := portaudio.OpenDefaultStream(1, 0, 88200, len(in), in)
	check(err)
	defer stream.Close()

	check(stream.Start())
	f := 800 / float32(len(in))
	amplit := 500000000
	g := 600 / float32(amplit)
	for {
		check(stream.Read())

		// Render values in window
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		renderer.SetDrawColor(255, 0, 0, 255)
		for i := range in {
			x := int(float32(i) * f)
			y := int32(float32(in[x])*g + 300)
			renderer.DrawPoint(int32(x), y)
		}
		renderer.Present()

		// See https://stackoverflow.com/questions/39637824/border-titlebar-not-properly-displaying-in-sdl-osx
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
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
