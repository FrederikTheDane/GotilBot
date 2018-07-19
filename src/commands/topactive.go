package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"helpers"
	"time"
)

type TopActive struct {
	running  bool
	listener chan []string
	msgChan chan map[string]int
	guildID  string
	tChannels helpers.ChannelArr
}

func (t *TopActive) Help () string {
	return "Use 'topactive <start|stop>' to enable or disable top active ranking"
}

func (t *TopActive) Run (uMember helpers.UMember, message *discordgo.Message, args []string) {
	var channels helpers.ChannelArr
	t.SetGuild(uMember.Member.GuildID)
	t.listener = t.GetChannel()
	channels = uMember.GetGuildChannels()
	t.tChannels = channels.GetTextChannels()

	if t.msgChan == nil {
		t.msgChan = make(chan map[string]int, 1)
	}

	if !t.running {
		go t.RunPassive(uMember, message)
	}

	t.listener <- args
	return
}

func (t *TopActive) Name () string {
	return "topactive"
}

func (t *TopActive) UserPermissions () (int, int) {
	return discordgo.PermissionManageRoles, helpers.GuildPerms
}

func (t *TopActive) RunPassive(uMember helpers.UMember, message *discordgo.Message) {
	var maxKey string
	var maxInt int
	tempmap := make(map[string]int)
	top := make(map[string]int)
	t.running = true
	listener := t.GetChannel()
	weekAgo := time.Now().Add(-24*7*time.Hour)
	go t.GetAllMessages(uMember, weekAgo)
	for {
		select {
		case args := <- listener:
			arg := args[0]
			switch arg {
			case "start":
				t.running = true
			case "stop":
				t.running = false
			}
		default:
			if t.running {
				updated := <-t.msgChan
				for k, v := range updated{
					tempmap[k] = v
					fmt.Println(k, v)
				}
				for i := 0; i < 5; i++ {
					maxKey, maxInt = getMaxValue(tempmap)
					top[maxKey] = maxInt
					delete(tempmap, maxKey)
					fmt.Printf("\n%v : %v \n", maxKey, maxInt)
				}
			} else {
				return
			}
		}
	}
}

func (t *TopActive) SetGuild(ID string) {
	t.guildID = ID
}

func (t *TopActive) GetGuildID() string {
	return t.guildID
}

func (t *TopActive) GetChannel() chan []string {
	if t.listener == nil {
		t.listener = make(chan []string)
	}
	return t.listener
}

func (t *TopActive) IsRunning() bool {
	return t.running
}

func (t *TopActive) GetAllMessages(uMember helpers.UMember, after time.Time) {
	var histories []helpers.ChannelHistory
	var totalMsg map[string]int
	channels := uMember.GetGuildChannels()
	textChannels := channels.GetTextChannels()
	totalMsg = make(map[string]int)

	for _, v := range textChannels {
		histories = append(histories,
			helpers.ChannelHistory{
			Session:      uMember.Session,
			MessageCount: make(map[string]int),
			ChannelID:    v.ID,
		})
	}
	for t.running{
		for _, v := range histories {
			//fmt.Printf("%T\n%v\n", v.MessageCount, v.ChannelID)
			v.InsertMessages(after)
			for k, v := range v.MessageCount {
				totalMsg[k] += v
			}
		}

		t.msgChan <- totalMsg

		/*
		for k := range totalMsg {
			totalMsg[k] = 0
		}
		*/

		totalMsg = make(map[string]int)
	}
}

func getMaxValue(m map[string]int) (string, int){
	var maxKey string
	var maxInt int
	for k, v := range m {
		if v > maxInt {
			maxKey = k
			maxInt = v
		}
	}
	return maxKey, maxInt
}