package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"helpers"
)

type TopActive struct {
	running bool
}

func (t *TopActive) Help () string {
	return "Use 'topactive <start|stop>' to enable or disable top active ranking"
}

func (t *TopActive) Run (uMember helpers.UMember, message *discordgo.Message, args []string) {
	for t.running {
		fmt.Println("Stopping Top Active ranking")
		return
	}
	return
}

func (t *TopActive) Name () string {
	return "topactive"
}

func (t *TopActive) UserPermissions () int {
	return discordgo.PermissionManageRoles
}

func GetAllMessages(channelID string, timestamp discordgo.Timestamp) {

}