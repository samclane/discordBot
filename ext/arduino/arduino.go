package arduino

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/tarm/serial"
	"log"
)

const (
	DISCONNECTED = iota
	CONNECTED    = iota
	MUTED        = iota
	DEAFENED     = iota
)

type statusCode byte

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
		return nil, err
	}
	return s, nil
}

func (a *Arduino) OnVoiceStateUpdate(_ *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	vs := vsu.VoiceState

	var sc statusCode

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
	fmt.Println([]byte(string(byte(sc))))
	_, err = s.Write([]byte(string(byte(sc))))
	if err != nil {
		log.Fatal(err)
	}

	err = s.Flush()
	if err != nil {
		log.Fatal(err)
	}

	err = s.Close()
	if err != nil {
		log.Fatal(err)
	}
}
