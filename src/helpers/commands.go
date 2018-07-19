package helpers

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
)

type Command interface {
	Help() string
	Run(uMember UMember, message *discordgo.Message, args []string)
	Name() string
	UserPermissions() (int, int)
}

type PassiveCommand interface {
	Command
	RunPassive(uMember UMember, message *discordgo.Message)
	SetGuild(ID string)
	GetGuildID() string
	GetChannel() chan []string
	IsRunning() bool
}


func Invoke (session *discordgo.Session, command Command, uMember UMember, message *discordgo.Message, args []string) chan []string {
	fmt.Println("Trying to invoke the command")
	channel, ok := session.State.Channel(message.ChannelID)
	if ok != nil {
		return nil
	}
	reqPerms, scope := command.UserPermissions()

	_, isPassive := command.(PassiveCommand)

	if uMember.HasPerms(reqPerms, scope, channel) {
		fmt.Println("Got the perms!")
		if isPassive {
			fmt.Println("It's passive!")
			command := command.(PassiveCommand)
			command = uMember.Guild.RunningCommands[command.Name()]

			fmt.Printf("%T\n", uMember.Guild.CommandsChannels[command.Name()])
			fmt.Printf("%T\n", command.GetChannel())

			command.SetGuild(uMember.Guild.Guild.ID)
			go command.Run(uMember, message, args)
			return uMember.Guild.CommandsChannels[command.Name()]
		}
		go command.Run(uMember, message, args)
		return nil
	}
	return nil
}

