package commands

import (
	"github.com/bwmarrin/discordgo"
	"helpers"
)


type Command interface {
	Help() string
	Run(uMember helpers.UMember, message *discordgo.Message, args []string)
	Name() string
	UserPermissions() int
}

func Invoke (session *discordgo.Session, command Command, uMember helpers.UMember, message *discordgo.Message, args []string) {
	channel, ok := session.State.Channel(message.ChannelID)
	if ok != nil {
		return
	}
	if uMember.HasPerms(command.UserPermissions(), channel) {
		command.Run(uMember, message, args)
	}
}