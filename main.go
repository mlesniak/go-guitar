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
	check(err)
	defer stream.Close()

	sampleN := 0
	//f, _ := os.Create("out.csv")

	check(stream.Start())
	for {
		check(stream.Read())

		buffer := new(bytes.Buffer)
		for i, v := range in {
			pos := sampleN*len(in) + i
			buffer.WriteString(fmt.Sprintf("%d,%d\n", pos, v))
		}
		fmt.Println(buffer.String())

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

func check(err error) {
	if err != nil {
		panic(err)
	}
}
