package main

import (
	"github.com/samclane/drivers/max7219"
)

func main() {
	ma := max7219.New(2, 3, 4)
	ma.MaxSingle(0x01, 0b10000001)
	ma.MaxSingle(0x02, 0b01000010)
	ma.MaxSingle(0x03, 0b00100100)
	ma.MaxSingle(0x04, 0b00011000)
	ma.MaxSingle(0x05, 0b00011000)
	ma.MaxSingle(0x06, 0b00100100)
	ma.MaxSingle(0x07, 0b01000010)
	ma.MaxSingle(0x08, 0b10000001)

}
