package bot

import (
	"github.com/bwmarrin/discordgo"
	"commands"
)

type Bot struct {
	Session *discordgo.Session
	User *discordgo.User
	Prefix string
	Commands map[string]commands.Command
}