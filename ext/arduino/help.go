package arduino

import (
	"discordBot/ext/mux"
	"github.com/bwmarrin/discordgo"
	"log"
)

func (a *SerialInterface) Help(ds *discordgo.Session, dm *discordgo.Message, _ *mux.Context) {
	_, err := ds.ChannelMessageSend(dm.ChannelID, "Sends voice state updates to an Arduino")
	if err != nil {
		log.Fatal(err)
	}
	return
}
