package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"os"
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
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

	renderer.Clear()
	renderer.SetDrawColor(255, 0, 0, 255)

	arr := make([]int32, 800)
	arr[0] = 300
	for i, _ := range arr {
		if i > 0 && i < len(arr)-2 {
			if rand.Int()%2 == 0 {
				arr[i] = arr[i-1] + 1
			} else {
				arr[i] = arr[i-1] - 1
			}
		}
	}

	for i, _ := range arr {
		renderer.DrawPoint(int32(i), arr[i])
	}

	fmt.Println("Renderer presenting")
	renderer.Present()
	fmt.Println("Renderer Waiting")

	// Remembered for later.
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}
