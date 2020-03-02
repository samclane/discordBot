package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const Version = "v0.0.0-alpha"

var Session, _ = discordgo.New()

func init() {
	// Discord Auth Token
	Session.Token = os.Getenv("DS_TOKEN")
	if Session.Token == "" {
		flag.StringVar(&Session.Token, "t", "", "Discord Authentication Token")
	}
}

func main() {
	var err error

	// Print fancy logo
	fmt.Printf(` 
	 _______               __      _______               __     
	/       \             /  |    /       \             /  |    
	$$$$$$$  |  ______   _$$ |_   $$$$$$$  |  ______   _$$ |_   
	$$ |__$$ | /      \ / $$   |  $$ |__$$ | /      \ / $$   |  
	$$    $$<  $$$$$$  |$$$$$$/   $$    $$< /$$$$$$  |$$$$$$/   
	$$$$$$$  | /    $$ |  $$ | __ $$$$$$$  |$$ |  $$ |  $$ | __ 
	$$ |  $$ |/$$$$$$$ |  $$ |/  |$$ |__$$ |$$ \__$$ |  $$ |/  |
	$$ |  $$ |$$    $$ |  $$  $$/ $$    $$/ $$    $$/   $$  $$/ 
	$$/   $$/  $$$$$$$/    $$$$/  $$$$$$$/   $$$$$$/     $$$$/  
	%-16s`+"\n\n", Version)

	// Parse command line arguments
	flag.Parse()

	// Verify a Token was provided
	if Session.Token == "" {
		log.Println("You must provide a Discord authentication token.")
		return
	}

	// Open a websocket connection to Discord
	err = Session.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		os.Exit(1)
	}

	// Wait for a CTRL-C
	log.Printf(`Now running. Press CTRL-C to exit.`)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-signals

	log.Printf(`Exit signal received. Cleaning up.`)
	// Clean up
	err = Session.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(`Closed.`)
}
