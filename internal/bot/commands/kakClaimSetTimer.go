package commands

import (
	"fmt"
	"slices"
	"time"

	"github.com/Sush1sui/fns-go/internal/common"
	"github.com/bwmarrin/discordgo"
)

func KakClaimSetTimer(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		return
	}


	isAuthorized := false
	for _, role := range i.Member.Roles {
		if slices.Contains(common.StaffRoleIDs, role) {
			isAuthorized = true
			break
		}
	}
	if !common.IsGuildOwner(&discordgo.Guild{ID: i.GuildID}, i.Member.User.ID) {
		isAuthorized = true
	}

	if !isAuthorized {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You do not have permission to set the kak claim timer.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	newTimer := i.ApplicationCommandData().GetOption("timer")

	common.KakClaimTimer = time.Duration(newTimer.IntValue()) * time.Second


	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Kak claim timer set to " + fmt.Sprintf("%d", newTimer.IntValue()) + " seconds.",
		},
	})

	if err != nil {
		fmt.Println("Error responding to interaction:", err)
		return
	}
}