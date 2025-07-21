package commands

import (
	"fmt"

	"github.com/Sush1sui/fns-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func StickyDelete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		return
	}

	channel := i.ApplicationCommandData().GetOption("channel")
	if channel == nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You must specify a channel to delete from sticky channels.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Printf("Error responding to interaction: %v\n", err)
		}
		return
	}

	count, e := repository.StickyService.DBClient.DeleteStickyChannel(channel.ChannelValue(s).ID)
	if e != nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to delete sticky channel.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Printf("Error responding to interaction: %v\n", err)
		}
		return
	}

	content := "Sticky channel deleted successfully\nCount: " + fmt.Sprintf("%d", count)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		fmt.Printf("Error responding to interaction: %v\n", err)
		return
	}
}