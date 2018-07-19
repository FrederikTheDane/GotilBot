package commands

import (
	"github.com/bwmarrin/discordgo"
	"helpers"
	"fmt"
	"time"
	"strconv"
)

type TestPerms struct {}

func (t *TestPerms) Help() string {
	return "Command for testing perms"
}

func (t *TestPerms) Run(uMember helpers.UMember, m *discordgo.Message, args []string) {
	uMember.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
		"```" +
		"You have the base perms: \n" +
		"%b \n" +
		"\n" +
		"You have the listener perms: \n" +
		"%b \n" +
		"\n" +
		"Total perms: \n" +
		"%b \n" +
		"```",
		uMember.GetBasePerms(), uMember.GetChannelPerms(m.ChannelID), uMember.GetBasePerms() | uMember.GetChannelPerms(m.ChannelID)))
}

func (t *TestPerms) Name() string {
	return "testperms"
}

func (t *TestPerms) UserPermissions() (int, int) {
	return 0, helpers.BasePerms
}



type ChannelMessages struct {}

func (*ChannelMessages) Help() string {
	return "Run this command to test the TopActive message getter"
}

func (*ChannelMessages) Run(uMember helpers.UMember, message *discordgo.Message, args []string) {
	if len(args) == 0 {
		uMember.Session.ChannelMessageSend(message.ChannelID, "Please provide a number")
		return
	}
	limitstring := args[0]
	limit, err := strconv.Atoi(limitstring)
	if err != nil || limit > 100 || limit < 1{
		uMember.Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%v is not a valid number! Choose a number from 1 to 100", limitstring))
		return
	}
	msgs, err := uMember.Session.ChannelMessages(message.ChannelID, limit, "", "", "")
	var msgsTime []string
	for k, v := range msgs {
		msgsTime = append(msgsTime, fmt.Sprintf("%v: %v \n", k, v.Timestamp))
	}
	uMember.Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%v, %v", msgsTime, err))
	return
}

func (*ChannelMessages) Name() string {
	return "chanmessages"
}

func (*ChannelMessages) UserPermissions() (int, int) {
	return discordgo.PermissionAdministrator, helpers.GuildPerms
}



type TimeTest struct {

}

func (*TimeTest) Help() string {
	return "Some simple tests for Discord Timestamps"
}

func (*TimeTest) Run(uMember helpers.UMember, message *discordgo.Message, args []string) {
	uMember.Session.ChannelMessageSend(message.ChannelID, "Fetching activity... This will take a while")
	weekAgo := time.Now().Add(-24*7*time.Hour)
	uMember.Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("A week ago the time was %v", weekAgo))
	msgMap := make(map[string]int)
	history := helpers.ChannelHistory{
		Session:      uMember.Session,
		MessageCount: msgMap,
		ChannelID:    message.ChannelID,
	}
	uMember.Session.ChannelMessageSend(message.ChannelID, "Mapping messages to users... This is the heavy part")
	history.InsertMessages(weekAgo)
	goTime, _ := message.Timestamp.Parse()
	uMember.Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("```" +
		"Your message timestamp: \n" +
		"%#v \n" +
		"\n" +
		"Your message timestamp, fomatted as a Go timestamp \n" +
		"%#v \n" +
		"\n" +
		"Messages per user this past week:" +
		"```", message.Timestamp, goTime))
	for k, v := range history.MessageCount {
		uMember.Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("`%v: %d`", k, v))
	}
}

func (*TimeTest) Name() string {
	return "timetest"
}

func (*TimeTest) UserPermissions() (int, int) {
	return discordgo.PermissionAdministrator, helpers.GuildPerms
}
