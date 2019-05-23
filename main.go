package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cabrinha/commands/ping"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

func init() {
	// Setup our config file and read it
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error reading config file: %s", err))
	}
}

func main() {
	goBot, err := discordgo.New("Bot " + viper.GetString("token"))
	if err != nil {
		fmt.Println("error creating discord session: ", err)
	}

	goBot.AddHandler(messageCreate)

	err = goBot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	router := exrouter.New()
	router.On("ping", ping.PingRoute)
	router.On("pong", ping.PongRoute)

	goBot.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		router.FindAndExecute(goBot, viper.GetString("prefix"), goBot.State.User.ID, m.Message)
	})

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	goBot.Close()

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
}
