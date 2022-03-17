package utils

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v2"
)

func SendReply(s *discordgo.Session, m *discordgo.MessageCreate, msg string) {
	s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
}

func SendEmbed(s *discordgo.Session, m *discordgo.MessageCreate, title, msg string) {
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embeds:    []*discordgo.MessageEmbed{embed.NewGenericEmbed(title, msg)},
		Reference: m.Reference(),
	})
}

func GetData(key string) []string {
	output := []string{}
	file, err := ioutil.ReadFile("assets/" + key + ".yaml")
	if err != nil {
		log.Errorf("error: %v", err)
		return []string{}
	}
	yaml.Unmarshal(file, &output)
	return output
}

func RandomNumberBetween(n, m int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(m-n) + n
}

func RandomChoice(items []string) string {
	rand.Seed(time.Now().Unix())
	return items[rand.Intn(len(items))]
}

func YesNo() bool {
	rand.Seed(time.Now().Unix())
	return rand.Intn(2) == 1
}

func GambleMoney(amount int) (bool, int) {
	if RandomNumberBetween(0, 10) > 6 {
		loss := RandomNumberBetween(0, amount)
		return false, amount - loss
	}
	win := RandomNumberBetween(0, amount)
	return true, amount + win
}

func ParseAmount(m *discordgo.MessageCreate, pos int) (int, bool, error) {
	msgParts := strings.Split(m.Content, " ")
	if len(msgParts) > pos {
		amount, err := strconv.Atoi(msgParts[pos])
		if err != nil {
			return 0, true, errors.New("invalid amount")
		}
		if amount < 0 {
			return 0, true, errors.New("invalid amount")
		}
		return amount, true, nil
	}
	return 0, false, nil
}
