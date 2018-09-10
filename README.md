# Overview

## Goal

Write a go program to distinguish between three different guitar chords (Em, Am, G) and recognize them.

## Installation

    brew install portaudio
    go get
    go build

    ./go-guitar demo.aiff
    ... time passes
    afplay demo.aiff


## Open questions

- what is sampling actually?
- what does 44100 bits/second actually mean?


## References

- [AIFF documentation](http://www-mmsp.ece.mcgill.ca/Documents/AudioFormats/AIFF/Docs/AIFF-1.3.pdf)
