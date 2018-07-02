package commands

import (
	"github.com/bwmarrin/discordgo"
	"helpers"
)

type Ping struct {}

func (p *Ping) Help() string {
	return "Test if the bot is running"
}

func (p *Ping) Run(uMember helpers.UMember, message *discordgo.Message, args []string) {
	uMember.Session.ChannelMessageSend(message.ChannelID, "Pong!")
}

func (p *Ping) Name() string {
	return  "ping"
}

func (p *Ping) UserPermissions() int {
	return 0
}