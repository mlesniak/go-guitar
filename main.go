package main

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
	"os"
	"os/signal"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	fmt.Println("Starting recording, press CTRL+C to abort.")

	portaudio.Initialize()
	defer portaudio.Terminate()

	in := make([]int32, 64)
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	chk(err)
	defer stream.Close()

	chk(stream.Start())
	for {
		chk(stream.Read())
		// TODO ML Display data
		select {
		case <-sig:
			return
		default:
		}
	}
	chk(stream.Stop())
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
