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

	in := make([]int32, 1024)
	stream, err := portaudio.OpenDefaultStream(1, 0, 1024, len(in), in)
	chk(err)
	defer stream.Close()

	chk(stream.Start())
	for {
		chk(stream.Read())
		// TODO ML Display data
		min, max := stats(in)
		fmt.Println("SAMPLE:", min, max)

		// Check if we should exit?
		select {
		case <-sig:
			return
		default:
		}
	}
	chk(stream.Stop())
}

func stats(arr []int32) (int32, int32) {
	min := arr[0]
	for _, v := range arr {
		if v < min {
			min = v
		}
	}

	max := arr[0]
	for _, v := range arr {
		if v > max {
			max = v
		}
	}

	return min, max
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
