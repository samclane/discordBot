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

const (
	DISCONNECTED = iota
	CONNECTED    = iota
	MUTED        = iota
	DEAFENED     = iota
)

type Arduino struct {
	comPort      string
	baudRate     int
	serialConfig *serial.Config
}

func New(cp string, br int) *Arduino {
	a := &Arduino{
		baudRate: br,
		comPort:  cp,
	}
	a.serialConfig = &serial.Config{
		Name: cp,
		Baud: br,
	}

	return a
}

func (a *Arduino) SerialConnect() (*serial.Port, error) {
	s, err := serial.OpenPort(a.serialConfig)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(2 * time.Second)
	return s, nil
}

func (a *Arduino) OnVoiceStateUpdate(_ *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
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

	// TODO Make this send bytes to Arduino correctly
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
