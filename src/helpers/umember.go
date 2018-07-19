package helpers

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
)

type UMember struct {
	Session *discordgo.Session
	Guild *UGuild
	Member *discordgo.Member
	ChannelPerms map[*discordgo.Channel]int
}

func (uMember *UMember) HasPerms(perms int, scope int, channel *discordgo.Channel) bool {
	memPerms := uMember.GetTotalChannelPerms(channel.ID)
	if perms == 0 || memPerms & discordgo.PermissionAdministrator > 0 {
		return true
	}
	switch scope {
	case BasePerms:
		return uMember.GetBasePerms() & perms == perms
	case ChannelPerms:
		return uMember.GetChannelPerms(channel.ID) & perms == perms
	case TotalChannelPerms:
		return uMember.GetTotalChannelPerms(channel.ID) & perms == perms
	case GuildPerms:
		return uMember.GetTotalGuildPerms() & perms == perms
	}

	return false
}

func (uMember *UMember) GetTotalGuildPerms() int {
	var perms int
	for _, v := range uMember.GetGuildChannels() {
		perms |= uMember.GetChannelPerms(v.ID)
	}
	perms |= uMember.GetBasePerms()
	return perms
}

func (uMember *UMember) GetTotalChannelPerms(channel string) int {
	return uMember.GetBasePerms() | uMember.GetChannelPerms(channel)
}

func (uMember *UMember) GetChannelPerms (channel string) int {
	perms, _ := uMember.Session.State.UserChannelPermissions(uMember.Member.User.ID, channel)
	return perms
}

func (uMember *UMember) GetBasePerms() int {
	var perms int
	for _, v := range uMember.GetRoleList() {
		perms |= v.Permissions
	}
	return perms
}

func (uMember *UMember) GetRoleList() discordgo.Roles{
	var roles discordgo.Roles
	roleIDS := uMember.Member.Roles
	guildID := uMember.Member.GuildID
	for _, v := range roleIDS {
		role, _ := uMember.Session.State.Role(guildID, v)
		roles = append(roles, role)
	}
	return roles
}

func (uMember *UMember) GetGuildChannels() ChannelArr {
	channels, ok := uMember.Session.GuildChannels(uMember.Member.GuildID)
	if ok != nil {
		fmt.Println(ok)
		return []*discordgo.Channel{}
	}
	return channels
}

func GetUMember (s *discordgo.Session, m discordgo.Message, uGuild *UGuild) *UMember {
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
	uMember := &UMember{
		Session:      s,
		Guild:        uGuild,
		Member:       member,
		ChannelPerms: channelPerms,
	}
	return uMember
}