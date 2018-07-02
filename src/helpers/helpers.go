package helpers

import (
	"github.com/bwmarrin/discordgo"
	"time"
	"fmt"
)

type UMember struct {
	Session *discordgo.Session
	Member *discordgo.Member
	ChannelPerms map[*discordgo.Channel]int
}

func (uMember UMember) HasPerms(perms int, channel *discordgo.Channel) bool {
	memPerms := uMember.GetTotalPerms(channel.ID)
	if perms == 0 || memPerms & discordgo.PermissionAdministrator > 0 {
		return true
	}
	return memPerms & perms == perms
}

func (uMember UMember) GetTotalPerms (channel string) int {
	return uMember.GetBasePerms() | uMember.GetChannelPerms(channel)
}

func (uMember UMember) GetChannelPerms (channel string) int {
	perms, _ := uMember.Session.State.UserChannelPermissions(uMember.Member.User.ID, channel)
	return perms
}

func (uMember UMember) GetBasePerms() int {
	var perms int
	for _, v := range uMember.GetRoleList() {
		perms |= v.Permissions
	}
	return perms
}

func (uMember UMember) GetRoleList() discordgo.Roles{
	var roles discordgo.Roles
	roleIDS := uMember.Member.Roles
	guildID := uMember.Member.GuildID
	for _, v := range roleIDS {
		role, _ := uMember.Session.State.Role(guildID, v)
		roles = append(roles, role)
	}
	return roles
}

func GetGuild (s *discordgo.Session, channelID string) discordgo.Guild {
	channel, _ := s.Channel(channelID)
	guild, _ := s.Guild(channel.GuildID)
	return *guild
}

func GetUMember (s *discordgo.Session, m discordgo.Message) UMember {
	state := s.State
	guild := GetGuild(s, m.ChannelID)
	member, _ := s.State.Member(guild.ID, m.Author.ID)
	channelPerms := make(map[*discordgo.Channel]int)
	for _, v := range guild.Channels{
		perms, success := state.UserChannelPermissions(m.Author.ID, m.ChannelID)
		if success != nil {
			panic(success)
		}
		channelPerms[v] = perms
	}
	uMember := UMember{
		Session:      s,
		Member:       member,
		ChannelPerms: channelPerms,
	}
	return uMember
}

type ChannelHistory struct {
	Session *discordgo.Session
	MessageCount map[string]int
	ChannelID string
}


func (c *ChannelHistory) InsertMessages (oldest time.Time) {
	var messArr []*discordgo.Message
	var timestamp time.Time
	var err error
	var prevmsg *discordgo.Message
	var prev []*discordgo.Message
	c.MessageCount = make(map[string]int)

	timestamp = time.Now()
	prev, _ = c.Session.ChannelMessages(c.ChannelID, 1, "", "", "")
	prevmsg = prev[0]

	for timestamp.After(oldest) {
		messArr = append(messArr, prev[0])

		prev, _ = c.Session.ChannelMessages(c.ChannelID, 1, prevmsg.ID, "", "")
		prevmsg = prev[0]

		if prev[0] == nil {
			break
		}
		timestamp, err = prev[0].Timestamp.Parse()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%v \n%v \n%v \n", len(messArr), timestamp, messArr[len(messArr)-1:][0].ID)
	}

	for _, v := range messArr {
		c.MessageCount[v.Author.String()] += 1
	}
}
