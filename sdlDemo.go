package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"os"
	"time"
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

	go func() {
		for true {
			renderer.SetDrawColor(0, 0, 0, 255)
			renderer.Clear()
			renderer.SetDrawColor(255, 0, 0, 255)
			arr := make([]int32, 8192)
			arr[0] = 300
			for i := range arr {
				if i > 0 && i < len(arr)-2 {
					if rand.Int()%2 == 0 {
						arr[i] = arr[i-1] + 1
					} else {
						arr[i] = arr[i-1] - 1
					}
				}
			}

			f := 800 / float32(8192)
			for i := range arr {
				// Compute index
				// 0 -> 0
				// x -> 800/8192
				// 8192 -> 800

				x := int(float32(i) * f)
				renderer.DrawPoint(int32(x), arr[x])
			}
			renderer.Present()
			time.Sleep(100 * time.Millisecond)
		}
	}()

	fmt.Println("Renderer presenting")

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
