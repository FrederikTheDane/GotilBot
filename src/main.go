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
	dg.AddHandler(onGuildJoin)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	botCommands := make(map[string]helpers.Command)
	botPassiveCommands := make(map[string]helpers.PassiveCommand)

	passiveCommandList := []helpers.PassiveCommand{
		&commands.TopActive{},
	}

	for _, v := range passiveCommandList {
		botPassiveCommands[v.Name()] = v
	}

	commandList := []helpers.Command{
		&commands.TestPerms{},
		&commands.Ping{},
		&commands.TimeTest{},
		&commands.ChannelMessages{},
		&commands.RussianRoulette{},
	}

	for _, v := range commandList {
		botCommands[v.Name()] = v
	}


	guildsMap := make(map[string]*helpers.UGuild)

	for _, v := range dg.State.Guilds {
		uGuild := helpers.GetUGuild(v)
		guildsMap[v.ID] = uGuild
		for _, c := range passiveCommandList {
			guildsMap[v.ID].RunningCommands[c.Name()] = c
			guildsMap[v.ID].CommandsChannels[c.Name()] = make(chan []string)
		}
	}

	discordBot = &bot.Bot{
		Session:         dg,
		User:            dg.State.User,
		Prefix:          Prefix,
		Commands:        botCommands,
		PassiveCommands: botPassiveCommands,
		Guilds:          guildsMap,
	}

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
	if len(m.Content) <= len(Prefix){
		return
	}

	if m.Content[:len(Prefix)] == Prefix {
		if err :=  icommands(s, message); err != nil {
			fmt.Println(err)
		}
	}
}

func onGuildJoin(s *discordgo.Session, g *discordgo.GuildCreate) {

}

func icommands(s *discordgo.Session, m discordgo.Message) error {
	var command helpers.Command
	var exists bool

	uMember := helpers.GetUMember(s, m, discordBot.Guilds[helpers.GetGuild(s, m.ChannelID).ID])

	msgFields := strings.Fields(m.Content[len(Prefix):])
	args := msgFields[1:]



	if command, exists = discordBot.Commands[msgFields[0]]; exists {
		helpers.Invoke(s, command, *uMember, &m, args)
		return nil
	} else if command, exists = discordBot.PassiveCommands[msgFields[0]]; exists {
		fmt.Println("Command exists!")
		guild := helpers.GetGuild(s, m.ChannelID)
		//uGuild := helpers.GetUGuild(guild)

		command := command.(helpers.PassiveCommand)
		command = discordBot.Guilds[guild.ID].RunningCommands[command.Name()]
		channel := discordBot.Guilds[guild.ID].CommandsChannels[command.Name()]
		channel = helpers.Invoke(s, command, *uMember, &m, args)
		channel <- args

		return nil
	}


	if !exists{
		return &CommandNotFoundError{
			time.Now(),
			m.Content[len(Prefix):],
			m.Author,
		}
	}
	return nil
}
