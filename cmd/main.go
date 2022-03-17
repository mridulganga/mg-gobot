package main

import (
	"fmt"

	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/mridulganga/mg-gobot/pkg/constants"
	"github.com/mridulganga/mg-gobot/pkg/db"
	"github.com/mridulganga/mg-gobot/pkg/handlers"
	"github.com/mridulganga/mg-gobot/pkg/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}
	log.SetLevel(log.DebugLevel)

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		log.Errorf("error creating Discord session,", err)
		return
	}

	monoHandler := handlers.NewMono(db.NewDB(constants.DBPath))

	router := router.NewRouter(handlers.NotFoundHandler)
	router.DirectRoute("hi,hello,sup,hey", handlers.GreetingsHandler)
	router.Route("help", handlers.HelpHandler)
	router.Route("balance", monoHandler.BalanceHandler)
	router.Route("rich", monoHandler.RichHandler)
	router.Route("send", monoHandler.SendHandler)
	router.Route("beg", monoHandler.BegHandler)
	router.Route("search", monoHandler.SearchHandler)
	router.Route("steal", monoHandler.StealHandler)
	router.Route("deposit", monoHandler.DepositHandler)
	router.Route("withdraw", monoHandler.WithdrawHandler)
	router.Route("gamble", monoHandler.GambleHandler)

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if !m.Author.Bot {
			log.Infof("%s: %s", m.Author.Username, m.Content)
			router.Resolve(s, m)
		}
	})

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Errorf("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
