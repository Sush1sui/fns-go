package commands

import (
	"fmt"

	"github.com/Sush1sui/fns-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func StickyDeleteAll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		return
	}

	count, err := repository.StickyService.DBClient.DeleteAllStickyChannels()
	if err != nil {
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to delete sticky channels.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Printf("Error responding to interaction: %v\n", e)
		}
		return
	}

	// Respond to the interaction
	e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("All sticky channels have been deleted successfully.\nCount: %d channels.", count),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if e != nil {
		fmt.Printf("Error responding to interaction: %v\n", e)
	}
}