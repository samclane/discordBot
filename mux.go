package main

import (
	"discordBot/ext/mux"
	"log"
)

// This file adds the Disgord message route multiplexer, aka "command router".
// to the Disgord bot. This is an optional addition however it is included
// by default to demonstrate how to extend the Disgord bot.

var Router = mux.New()

func init() {
	Session.AddHandler(Router.OnMessageCreate)
	_, err := Router.Route("help", "Display this message", Router.Help)
	if err != nil {
		log.Fatal(err)
	}
}
