package arduino

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/tarm/serial"
	"log"
	"time"
)

// Discord Statuses
const (
	DISCONNECTED = iota
	CONNECTED    = iota
	MUTED        = iota
	DEAFENED     = iota
)

// Arduino object definition
type SerialInterface struct {
	comPort      string         // Name of serial port Ex "dev/ttyACM0", "COM3"
	baudRate     int            // 9600 is a good starting value
	serialConfig *serial.Config // reuses the com-port and baud rate values
}

func New(cp string, br int) *SerialInterface {
	a := &SerialInterface{
		baudRate: br,
		comPort:  cp,
	}
	a.serialConfig = &serial.Config{
		Name: cp,
		Baud: br,
	}

	return a
}

func (a *SerialInterface) SerialConnect() (*serial.Port, error) {
	s, err := serial.OpenPort(a.serialConfig)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(2 * time.Second)
	return s, nil
}

func (a *SerialInterface) OnVoiceStateUpdate(_ *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	vs := vsu.VoiceState

	var sc int8

	if vs.Deaf || vs.SelfDeaf {
		sc = DEAFENED
	} else if vs.Mute || vs.SelfMute {
		sc = MUTED
	} else if vs.ChannelID != "" {
		sc = CONNECTED
	} else {
		sc = DISCONNECTED
	}

	s, err := a.SerialConnect()
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, sc); err != nil {
		log.Fatal(err)
	}
	n, err := s.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("sent %d byte(s): ", n)
	fmt.Println(buf.Bytes())
	err = s.Close()
	if err != nil {
		log.Fatal(err)
	}
}
