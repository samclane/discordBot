package main

import (
	"discordBot/ext/arduino"
	"log"
)

var ArduinoInterface = arduino.New("/dev/ttyACM0", 9600)

func init() {
	Session.AddHandler(ArduinoInterface.OnVoiceStateUpdate)
	_, err := Router.Route("ahelp", "Arduino Interface", ArduinoInterface.Help)
	if err != nil {
		log.Fatal(err)
	}
}
