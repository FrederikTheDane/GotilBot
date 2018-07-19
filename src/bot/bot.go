package bot

import (
	"github.com/bwmarrin/discordgo"
	"helpers"
)

type Bot struct {
	Session *discordgo.Session
	User *discordgo.User
	Prefix string
	Commands map[string]helpers.Command
	PassiveCommands map[string]helpers.PassiveCommand
	Guilds map[string]*helpers.UGuild
}
