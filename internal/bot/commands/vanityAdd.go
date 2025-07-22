package commands

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/Sush1sui/fns-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func VanityAdd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		return
	}
	
	user := i.ApplicationCommandData().GetOption("user").UserValue(s)
	if user == nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You must specify a user to add a vanity.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Failed to respond to interaction:", err)
			return
		}
		return
	}

	res, err := repository.ExemptedService.DBClient.ExemptUserVanity(user.ID)
	if err != nil || !res {
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to add vanity exemption",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Failed to respond to interaction:", e)
			return
		}
		return
	}

	newRoleIds := strings.Split(os.Getenv("SUPPORTER_ROLE_IDS"), ",")
	if len(newRoleIds) == 0 {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No supporter roles found in environment variables. Talk to your developer.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Failed to respond to interaction:", err)
			return
		}
		return
	}

	member, err := s.GuildMember(i.GuildID, user.ID)
	if err != nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to fetch user information.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Failed to respond to interaction:", err)
			return
		}
		return
	}

	updatedUserRoles := append(member.Roles, newRoleIds...)
	if !slices.Equal(member.Roles, updatedUserRoles) {
		_, err := s.GuildMemberEdit(i.GuildID, user.ID, &discordgo.GuildMemberParams{
			Roles: &updatedUserRoles,
		})
		if err != nil {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to update user roles.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				fmt.Println("Failed to respond to interaction:", err)
				return
			}
			return
		}
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Vanity added for user %s. They now have the supporter role.", user.Username),
			},
		})
		if err != nil {
			fmt.Println("Failed to respond to interaction:", err)
		}
		return
	}
}