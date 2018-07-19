package helpers

import (
	"github.com/bwmarrin/discordgo"
)

//Holds a map of pointers
type UGuild struct {
	Guild *discordgo.Guild
	RunningCommands map[string]PassiveCommand
	CommandsChannels map[string]chan []string
	Roles []*discordgo.Role
}

func (g *UGuild) RunPassive (command PassiveCommand) chan int {
	return nil
}

func GetUGuild(guild *discordgo.Guild) *UGuild {
	running := make(map[string]PassiveCommand)
	chans := make(map[string]chan []string)
	roles := guild.Roles
	uGuild := UGuild{
		Guild:            guild,
		RunningCommands:  running,
		CommandsChannels: chans,
		Roles:            roles,
	}
	return &uGuild
}