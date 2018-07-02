package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"bot"
	"commands"
	"helpers"
	"strings"
)

// Variables used for command line parameters
var (
	Token string
	Prefix string
	dg *discordgo.Session
	err error
	discordBot *bot.Bot
	botCommands map[string]commands.Command
)

type CommandNotFoundError struct {
	Timestamp time.Time
	Command string
	User *discordgo.User
}

func (e *CommandNotFoundError) Error() string {
	return fmt.Sprintf("At %v, the following command could not be found: %v. User: %v", e.Timestamp, e.Command, e.User)
}


func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&Prefix, "p", "", "Bot prefix")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err = discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}


	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(onMessage)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	botCommands = make(map[string]commands.Command)

	commandList := []commands.Command{
		&commands.TopActive{},
		&commands.TestPerms{},
		&commands.Ping{},
		&commands.TimeTest{},
	}

	for _, v := range commandList {
		botCommands[v.Name()] = v
	}


	discordBot = &bot.Bot{
		Session:  dg,
		User:     dg.State.User,
		Prefix:   Prefix,
		Commands: botCommands,
	}

	/*
	stop := make(chan bool)

	//go commands.TopActive(stop)

	time.Sleep(1*time.Second)

	stop <- true
	*/

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	message := *m.Message

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	//Process command, if the message is prefixed with the given Prefix
	if m.Content[:len(Prefix)] == Prefix {
		if err := icommands(s, message); err != nil {
			fmt.Println(err)
		}
	}
}


func icommands(s *discordgo.Session, m discordgo.Message) error {
	uMember := helpers.GetUMember(s, m)

	command, exists := discordBot.Commands[m.Content[len(Prefix):]]

	args := strings.Fields(m.Content)
	args = args[1:]

	if !exists{
		return &CommandNotFoundError{
			time.Now(),
			m.Content[len(Prefix):],
			m.Author,
		}
	}

	go commands.Invoke(s, command, uMember, &m, args)

	return nil
}
