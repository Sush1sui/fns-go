package commands

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/Sush1sui/fns-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func VanityRemove(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		return
	}

	user := i.ApplicationCommandData().GetOption("user").UserValue(s)
	if user == nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You must specify a user to remove a vanity.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Failed to respond to interaction:", err)
			return
		}
		return
	}

	count, err := repository.ExemptedService.DBClient.RemoveExemptedUser(user.ID)
	if err != nil || count == 0 {
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to remove vanity exemption",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Failed to respond to interaction:", e)
			return
		}
		return
	}

	roleIdsToRemove := strings.Split(os.Getenv("SUPPORTER_ROLE_IDS"), ",")
	if len(roleIdsToRemove) == 0 {
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

	supporterRoleIDs := strings.Split(os.Getenv("SUPPORTER_ROLE_IDS"), ",")
	if len(supporterRoleIDs) == 0 {
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
	rolesToRemove := map[string]struct{}{
		supporterRoleIDs[0]: {},
		supporterRoleIDs[1]: {},
	}
	updatedRoles := make([]string, 0, len(member.Roles))
	for _, roleId := range member.Roles {
		if _, remove := rolesToRemove[roleId]; !remove {
			updatedRoles = append(updatedRoles, roleId)
		}
	}

	if !slices.Equal(member.Roles, updatedRoles) {
		_, err := s.GuildMemberEdit(i.GuildID, user.ID, &discordgo.GuildMemberParams{
			Roles: &updatedRoles,
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
	} else {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("User %s does not have any vanity roles to remove.", user.Username),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Failed to respond to interaction:", err)
			return
		}
		return
	}
}