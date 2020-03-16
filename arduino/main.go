package main

import (
	"github.com/samclane/drivers/max7219"
	"time"
)

var X = [8]byte{
	0b10000001,
	0b01000010,
	0b00100100,
	0b00011000,
	0b00011000,
	0b00100100,
	0b01000010,
	0b10000001,
}

var PLAY = [8]byte{
	0b00000000,
	0b01100000,
	0b01111000,
	0b01111110,
	0b01111111,
	0b01111110,
	0b01111000,
	0b01100000,
}

var PAUSE = [8]byte{
	0b00000000,
	0b01100110,
	0b01100110,
	0b01100110,
	0b01100110,
	0b01100110,
	0b01100110,
	0b00000000,
}

var STOP = [8]byte{
	0b00000000,
	0b00000000,
	0b00111100,
	0b00111100,
	0b00111100,
	0b00111100,
	0b00000000,
	0b00000000,
}

func main() {
	ma := max7219.New(2, 3, 4)

	for i, row := range X {
		ma.MaxSingle(byte(i+1), row)
	}

	time.Sleep(time.Second * 1)

	for i, row := range PLAY {
		ma.MaxSingle(byte(i+1), row)
	}

	time.Sleep(time.Second * 1)

	for i, row := range PAUSE {
		ma.MaxSingle(byte(i+1), row)
	}

	time.Sleep(time.Second * 1)

	for i, row := range STOP {
		ma.MaxSingle(byte(i+1), row)
	}
}
