package main

import (
	"github.com/samclane/drivers/max7219"
	"github.com/tinygo-org/tinygo/src/machine"
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

var SMILE = [8]byte{
	0b00000000,
	0b00000000,
	0b01100110,
	0b01100110,
	0b00000000,
	0b01000010,
	0b00111100,
	0b00000000,
}

var (
	uart = machine.UART0
	tx   = machine.Pin(1) // 1, transmit line
	rx   = machine.Pin(0) // 0, receive line
)

func main() {
	ma := max7219.New(machine.Pin(2), machine.Pin(3), machine.Pin(4))
	ma.Configure()

	uart.Configure(machine.UARTConfig{
		BaudRate: 9600,
		TX:       tx,
		RX:       rx,
	})

	for {
		if uart.Buffered() > 0 {

			data, _ := uart.ReadByte() // read in a char

			switch data {
			case 0:
				for i, v := range X {
					ma.MaxSingle(byte(i+1), v)
				}
			case 1:
				for i, v := range PLAY {
					ma.MaxSingle(byte(i+1), v)
				}
			case 2:
				for i, v := range PAUSE {
					ma.MaxSingle(byte(i+1), v)
				}
			case 3:
				for i, v := range STOP {
					ma.MaxSingle(byte(i+1), v)
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
