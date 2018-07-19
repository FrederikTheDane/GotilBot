package commands

import (
	"helpers"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"time"
	"fmt"
)

type RussianRoulette struct {

}

func (*RussianRoulette) Help() string {
	return "Use this to gamble with your life"
}

func (*RussianRoulette) Run(uMember helpers.UMember, message *discordgo.Message, args []string) {
	chance := rand.New(rand.NewSource(time.Now().UnixNano()))
	outcome := chance.Intn(6)
	uMember.Session.ChannelMessageSend(message.ChannelID, fmt.Sprint(outcome))
	fmt.Println(outcome)
}

func (*RussianRoulette) Name() string {
	return "russianroulette"
}

func (*RussianRoulette) UserPermissions() (int, int) {
	return 0, helpers.BasePerms
}
