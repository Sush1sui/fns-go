package commands

import (
	"fmt"

	"github.com/Sush1sui/fns-go/internal/common"
	"github.com/bwmarrin/discordgo"
)

func StickyGetAll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		return
	}

	if len(common.StickyChannels) == 0 {
		embed := &discordgo.MessageEmbed{
			Title:       "No Sticky Channels Found",
			Description: "There are currently no sticky channels set up.",
			Color:       0xFF0000, // Red color
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
		return
	}

	var stickyList string
	count := 0
	for id := range common.StickyChannels {
		stickyList += "<#" + id + ">\n"
		count++
		if count < len(common.StickyChannels)-1 {
        stickyList += "\n"
    }
	}
	fmt.Print("Sticky Channels: ", stickyList)

	embed := &discordgo.MessageEmbed{
		Title:       "Sticky Channels",
		Description: stickyList,
		Color:       0x00FF00, // Green color
	}

	_, err := s.ChannelMessageSendEmbed(i.ChannelID, embed)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to fetch sticky channels.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Channels fetched successfully",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if e != nil {
		fmt.Printf("Error responding to interaction: %v\n", e)
		return
	}
}