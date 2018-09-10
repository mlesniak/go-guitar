package main

import (
	"bytes"
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

	n := 0
	chk(stream.Start())
	for {
		chk(stream.Read())
		// TODO ML Display data
		fmt.Println(toHex(in))
		n++

		// Check if we should exit?
		select {
		case <-sig:
			return
		default:
		}
	}
	chk(stream.Stop())
}

func toHex(arr []int32) string {
	var buffer bytes.Buffer
	for _, v := range arr {
		hex := fmt.Sprintf("%X ", v)
		buffer.WriteString(hex)
	}
	return buffer.String()
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}