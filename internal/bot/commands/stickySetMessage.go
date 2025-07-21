package commands

import (
	"fmt"

	"github.com/Sush1sui/fns-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func StickySetMessage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		return
	}


	channel := i.ApplicationCommandData().GetOption("channel")
	message := i.ApplicationCommandData().GetOption("message")

	if channel == nil || message == nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You must specify a channel and a message to create a sticky channel.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			s.ChannelMessageSend(i.ChannelID, "Failed to respond to interaction.")
			return
		}
		return
	}

	_, err := repository.StickyService.DBClient.SetStickyMessageId(channel.ChannelValue(s).ID, message.StringValue())
	if err != nil {
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to set sticky message",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Failed to respond to interaction:", e)
			return
		}
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Sticky message set successfully",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		fmt.Println("Failed to respond to interaction:", err)
		return
	}
}