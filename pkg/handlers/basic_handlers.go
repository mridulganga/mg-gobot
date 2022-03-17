package handlers

import (
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v2"
)

const (
	HelpFile = "assets/helps.yaml"
)

func NotFoundHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "I don't know that command yet.")
}

func GreetingsHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, m.Content)
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Reference: m.MessageReference,
	})
}

func HelpHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	helps := map[string]string{}
	helpFile, err := ioutil.ReadFile(HelpFile)
	if err != nil {
		log.Errorf("error: %v", err)
		return
	}
	yaml.Unmarshal(helpFile, &helps)

	mParts := strings.Split(m.Content, " ")

	if len(mParts) > 2 {
		if help, ok := helps[mParts[2]]; ok {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("MGBot Help", help))
			return
		}
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("MGBot Help", helps["basic"]))

}
