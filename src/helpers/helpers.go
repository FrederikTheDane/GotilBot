package helpers

import (
	"github.com/bwmarrin/discordgo"
	"time"
	"fmt"
)

func GetGuild (s *discordgo.Session, channelID string) *discordgo.Guild {
	channel, _ := s.Channel(channelID)
	guild, _ := s.Guild(channel.GuildID)
	return guild
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
	prev, _ = c.Session.ChannelMessages(c.ChannelID, 100, "", "", "")
	prevmsg = prev[len(prev)-1]

	for timestamp.After(oldest) {
		for _, v := range prev {
			tempStamp, _ := v.Timestamp.Parse()
			if tempStamp.Before(oldest) {
				break
			}
			messArr = append(messArr, v)
		}

		prev, _ = c.Session.ChannelMessages(c.ChannelID, 100, prevmsg.ID, "", "")
		if len(prev) == 0 {
			break
		}

		prevmsg = prev[len(prev)-1]

		timestamp, err = prev[0].Timestamp.Parse()
		if err != nil {
			fmt.Println(err)
		}

		//fmt.Printf("%v \n%v \n%v \n", len(messArr), timestamp, messArr[len(messArr)-1:][0].ID)
	}

	for _, v := range messArr {
		c.MessageCount[v.Author.String()] += 1
	}
	return
}

type ChannelArr []*discordgo.Channel

func (c *ChannelArr) GetTextChannels() ChannelArr {
	var textChannels ChannelArr
	for _, v := range *c {
		if v.Type == discordgo.ChannelTypeGuildText {
			textChannels = append(textChannels, v)
		}
	}
	return textChannels
}

const(
	BasePerms = iota
	ChannelPerms
	TotalChannelPerms
	GuildPerms
)