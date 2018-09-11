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

	// Initialize library.
	portaudio.Initialize()
	defer portaudio.Terminate()

	// Create buffer for sampling data and open the stream.
	in := make([]int32, 1024)
	stream, err := portaudio.OpenDefaultStream(1, 0, 8192, len(in), in)
	check(err)
	defer stream.Close()

	// Create file to store data.
	f, _ := os.Create("out.csv")
	defer f.Close()

	// Read sampleRate/len(in) samples per second.
	sampleN := 0
	check(stream.Start())
	for {
		check(stream.Read())

		// Create the CSV values for this sample and store them in file.
		buffer := new(bytes.Buffer)
		for i, v := range in {
			pos := sampleN*len(in) + i
			buffer.WriteString(fmt.Sprintf("%d,%d\n", pos, v))
		}
		f.WriteString(buffer.String())
		fmt.Println("#samples", sampleN)

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
