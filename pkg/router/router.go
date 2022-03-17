package router

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Router struct {
	routes       map[string]func(s *discordgo.Session, m *discordgo.MessageCreate)
	directRoutes map[string]func(s *discordgo.Session, m *discordgo.MessageCreate)
}

func NewRouter(notFoundHandler func(s *discordgo.Session, m *discordgo.MessageCreate)) Router {
	return Router{
		routes: map[string]func(s *discordgo.Session, m *discordgo.MessageCreate){
			"notfound": notFoundHandler,
		},
		directRoutes: map[string]func(s *discordgo.Session, m *discordgo.MessageCreate){},
	}
}

func (r Router) Route(startsWith string, f func(s *discordgo.Session, m *discordgo.MessageCreate)) {
	for _, sw := range strings.Split(startsWith, ",") {
		r.routes[sw] = f
	}
}

func (r Router) DirectRoute(startsWith string, f func(s *discordgo.Session, m *discordgo.MessageCreate)) {
	for _, sw := range strings.Split(startsWith, ",") {
		r.directRoutes[sw] = f
	}
}

func (r Router) Resolve(s *discordgo.Session, m *discordgo.MessageCreate) {
	mParts := strings.Split(m.Content, " ")
	if mParts[0] == "pls" {
		if f, ok := r.routes[mParts[1]]; ok {
			f(s, m)
			return
		}
		r.routes["notfound"](s, m)
	} else if f, ok := r.directRoutes[mParts[0]]; ok {
		f(s, m)
		return
	}
}
