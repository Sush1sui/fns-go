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

	// Defer the interaction so we can perform work and edit the response later
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		fmt.Println("Failed to defer interaction:", err)
		return
	}

	user := i.ApplicationCommandData().GetOption("user").UserValue(s)
	if user == nil {
		msg := "You must specify a user to remove a vanity."
		if _, e := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg}); e != nil {
			fmt.Println("Failed to edit deferred interaction:", e)
		}
		return
	}

	count, err := repository.ExemptedService.DBClient.RemoveExemptedUser(user.ID)
	if err != nil {
		msg := "Failed to remove vanity exemption"
		if _, e := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg}); e != nil {
			fmt.Println("Failed to edit deferred interaction:", e)
		}
		return
	}
	if count == 0 {
		// Nothing was removed — user had no exemption
		msg := fmt.Sprintf("User %s does not have an active vanity exemption.", user.Username)
		if _, e := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg}); e != nil {
			fmt.Println("Failed to edit deferred interaction:", e)
		}
		return
	}

	// Parse supporter role IDs once and validate
	supporterRoleIDs := strings.Split(os.Getenv("SUPPORTER_ROLE_IDS"), ",")
	if len(supporterRoleIDs) == 0 {
		msg := "No supporter roles found in environment variables. Talk to your developer."
		if _, e := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg}); e != nil {
			fmt.Println("Failed to edit deferred interaction:", e)
		}
		return
	}

	member, err := s.GuildMember(i.GuildID, user.ID)
	if err != nil {
		msg := "Failed to fetch user information."
		if _, e := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg}); e != nil {
			fmt.Println("Failed to edit deferred interaction:", e)
		}
		return
	}

	// Build rolesToRemove safely from available supporterRoleIDs (support 1 or more)
	rolesToRemove := make(map[string]struct{})
	if len(supporterRoleIDs) >= 1 && supporterRoleIDs[0] != "" {
		rolesToRemove[supporterRoleIDs[0]] = struct{}{}
	}
	if len(supporterRoleIDs) >= 2 && supporterRoleIDs[1] != "" {
		rolesToRemove[supporterRoleIDs[1]] = struct{}{}
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
			msg := "Failed to update user roles."
			if _, e := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg}); e != nil {
				fmt.Println("Failed to edit deferred interaction:", e)
			}
			return
		}

		// Edit deferred response with success after DB deletion and role update
		successMsg := fmt.Sprintf("✅ | Removed vanity exemption and roles for %s.", user.Username)
		if _, e := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &successMsg}); e != nil {
			fmt.Println("Failed to edit deferred interaction:", e)
		}
	} else {
		msg := fmt.Sprintf("User %s does not have any vanity roles to remove.", user.Username)
		if _, e := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg}); e != nil {
			fmt.Println("Failed to edit deferred interaction:", e)
		}
		return
	}
}