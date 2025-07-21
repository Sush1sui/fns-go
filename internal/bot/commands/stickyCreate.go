package commands

import (
	"fmt"

	"github.com/Sush1sui/fns-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func StickyCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		return
	}

	channel := i.ApplicationCommandData().GetOption("channel")
	message := i.ApplicationCommandData().GetOption("message")

	if channel == nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You must specify a channel to create a sticky channel.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			s.ChannelMessageSend(i.ChannelID, "Failed to respond to interaction.")
			return
		}
		return
	}

	// check if the channel is already a sticky channel
	channelID := channel.ChannelValue(s).ID
	isSticky := false
	for id := range repository.StickyChannels {
		if id == channelID {
			isSticky = true
			break
		}
	}
	if isSticky {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This channel is already a sticky channel.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			s.ChannelMessageSend(i.ChannelID, "Failed to respond to interaction.")
			return
		}
		return
	}

	if message == nil {
		_, err := repository.StickyService.DBClient.CreateStickyChannel(channelID, "")
		if err != nil {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to create sticky channel.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				fmt.Printf("Error responding to interaction: %v\n", err)
			}
			return
		}
	} else {
		_, err := repository.StickyService.DBClient.CreateStickyChannel(channelID, message.StringValue())
		if err != nil {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to create sticky channel.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				fmt.Printf("Error responding to interaction: %v\n", err)
			}
			return
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Sticky channel created successfully!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		fmt.Printf("Error responding to interaction: %v\n", err)
		return
	}
}