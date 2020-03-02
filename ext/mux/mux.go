// Package mux provides a simple Discord message route multiplexer that
// parses messages and then executes a matching registered handler, if found.
// mux can be used with both this bot and the DiscordGo library.
package mux

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"regexp"
	"strings"
)

// Route holds information about a specific message route handler
type Route struct {
	Pattern     string      // Match pattern that should trigger this route handler
	Description string      // short description of this route
	Help        string      // detailed help string for this route
	Run         HandlerFunc // route handler  function to call
}

// Context holds a bit of extra data we pass along to route handlers
// This way processing some of this only needs to happen once.
type Context struct {
	Fields          []string
	Content         string
	IsDirected      bool
	IsPrivate       bool
	HasPrefix       bool
	HasMention      bool
	HasMentionFirst bool
}

type HandlerFunc func(session *discordgo.Session, message *discordgo.Message, ctx *Context)

type Mux struct {
	Routes  []*Route
	Default *Route
	Prefix  string
}

// New returns a new Discord message route mux
func New() *Mux {
	m := &Mux{}
	m.Prefix = "-dg "
	return m
}

// Route allows you to register a route
func (m *Mux) Route(pattern, desc string, cb HandlerFunc) (*Route, error) {

	r := Route{}
	r.Pattern = pattern
	r.Description = desc
	r.Run = cb
	m.Routes = append(m.Routes, &r)

	return &r, nil
}

// FuzzyMatch attempts to find the best route match for a given message
func (m *Mux) FuzzyMatch(msg string) (*Route, []string) {

	// Tokenize the msg string into a slice of words
	fields := strings.Fields(msg)

	// stop if no fields
	if len(fields) == 0 {
		return nil, nil
	}

	// Search through the command list for a match
	var r *Route
	var rank int

	var fk int
	for fk, fv := range fields {
		for _, rv := range m.Routes {
			// If we find an exact match, return immediately
			if rv.Pattern == fv {
				return rv, fields[fk:]
			}

			// Some fuzzy logic searching
			if strings.HasPrefix(rv.Pattern, fv) {
				if len(fv) > rank {
					r = rv
					rank = len(fv)
				}
			}
		}
	}
	return r, fields[fk:]
}

func (m *Mux) OnMessageCreate(ds *discordgo.Session, mc *discordgo.MessageCreate) {

	var err error

	// Ignore bot
	if mc.Author.ID == ds.State.User.ID {
		return
	}

	ctx := &Context{
		Content: strings.TrimSpace(mc.Content),
	}

	var c *discordgo.Channel
	c, err = ds.State.Channel(mc.ChannelID)
	if err != nil {
		// Try fetching via REST API
		c, err = ds.Channel(mc.ChannelID)
		if err != nil {
			log.Printf("unable to fetch Channel for Message, %s", err)
		} else {
			// Attempt to add this channel into our State
			err = ds.State.ChannelAdd(c)
			if err != nil {
				log.Printf("error updating State with Channel, %s", err)
			}
		}
	}
	// Add Channel info into Context (if successfully get channel)
	if c != nil {
		if c.Type == discordgo.ChannelTypeDM {
			ctx.IsPrivate, ctx.IsDirected = true, true
		}
	}

	if !ctx.IsDirected {
		for _, v := range mc.Mentions {

			if v.ID == ds.State.User.ID {
				ctx.IsDirected, ctx.HasMention = true, true

				reg := regexp.MustCompile(fmt.Sprintf("<@!?(%s)>", ds.State.User.ID))

				// Was the @mention the first part of the string?
				if reg.FindStringIndex(ctx.Content)[0] == 0 {
					ctx.HasMentionFirst = true
				}

				// strip bot mention tags from content string
				ctx.Content = reg.ReplaceAllString(ctx.Content, "")

				break
			}
		}
	}

	// Detect prefix mention
	if !ctx.IsDirected && len(m.Prefix) > 0 {
		// TODO: No guild-defined prefix definition
		if strings.HasPrefix(ctx.Content, m.Prefix) {
			ctx.IsDirected, ctx.HasPrefix, ctx.HasMentionFirst = true, true, true
			ctx.Content = strings.TrimPrefix(ctx.Content, m.Prefix)
		}
	}

	// For now, if not specifically mentioned we do nothing
	if !ctx.IsDirected {
		return
	}

	// Try to find the best fuzzy-match command out of the message
	r, fl := m.FuzzyMatch(ctx.Content)
	if r != nil {
		ctx.Fields = fl
		r.Run(ds, mc.Message, ctx)
		return
	}

	// If not command match was found, call the default
	// Ignore if only @mentioned in the middle of message
	if m.Default != nil && (ctx.HasMentionFirst) {
		// TODO: This could be a rate-limit
		// Can cause an endless loop when talking to another bot
		m.Default.Run(ds, mc.Message, ctx)
	}

}
