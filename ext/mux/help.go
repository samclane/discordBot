package mux

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"sort"
	"strconv"
)

func (m *Mux) Help(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	// Set command prefix to display

	cp := ""
	if ctx.IsPrivate {
		cp = ""
	} else if ctx.HasPrefix {
		cp = m.Prefix
	} else {
		cp = fmt.Sprintf("@%s ", ds.State.User.Username)
	}

	// sort commands
	maxlen := 0
	keys := make([]string, 0, len(m.Routes))
	cmdmap := make(map[string]*Route)

	for _, v := range m.Routes {
		// Only display commands with a description
		if v.Description == "" {
			continue
		}

		// Calculate the max length of command+args string
		l := len(v.Pattern)
		if l > maxlen {
			maxlen = l
		}

		cmdmap[v.Pattern] = v

		// help and about are added separately
		if v.Pattern == "help" || v.Pattern == "about" {
			continue
		}

		keys = append(keys, v.Pattern)
	}

	sort.Strings(keys)

	// TODO: Learn more link needs to be configurable, currently points to original bot source
	resp := "\n*Commands can be abbreviated and mixed with other text. Learn more at <https://github.com/bwmarrin/disgord>*\n"
	resp += "```autoit\n"

	v, ok := cmdmap["help"]
	if ok {
		keys = append([]string{v.Pattern}, keys...)
	}

	v, ok = cmdmap["about"]
	if ok {
		keys = append([]string{v.Pattern}, keys...)
	}

	// Add sorted result to help msg
	for _, k := range keys {
		v := cmdmap[k]
		resp += fmt.Sprintf("%s%-"+strconv.Itoa(maxlen)+"s # %s\n", cp, v.Pattern+v.Help, v.Description)
	}

	resp += "```\n"

	_, err := ds.ChannelMessageSend(dm.ChannelID, resp)
	if err != nil {
		log.Fatal(err)
	}
	return
}
