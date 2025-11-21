package commands

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/Sush1sui/fns-go/internal/common"
	"github.com/Sush1sui/fns-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func VanityAdd(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	permanentOption := i.ApplicationCommandData().GetOption("permanent")
	permanent := false
	if permanentOption != nil {
		permanent = permanentOption.BoolValue()
	}
	
	user := i.ApplicationCommandData().GetOption("user").UserValue(s)
	if user == nil {
		msg := "You must specify a user to add a vanity."
		if _, e := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg}); e != nil {
			fmt.Println("Failed to edit deferred interaction:", e)
		}
		return
	}

	// Fetch the target member so we can check their actual roles.
	targetMember, err := s.GuildMember(i.GuildID, user.ID)
	if err != nil {
		msg := "Failed to fetch target user information."
		if _, e := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg}); e != nil {
			fmt.Println("Failed to edit deferred interaction:", e)
		}
		return
	}

	// Determine permanence based on the *target user's* staff membership.
	targetIsStaff := false
	for _, rid := range targetMember.Roles {
		if slices.Contains(common.StaffRoleIDs, rid) {
			targetIsStaff = true
			break
		}
		// If we reach here, the user already had the supporter role
		alreadyMsg := fmt.Sprintf("User <@%s> already has the supporter role or vanity is already active.", user.ID)
		if _, e := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &alreadyMsg}); e != nil {
			fmt.Println("Failed to edit deferred interaction:", e)
		}
	}

	if targetIsStaff || permanent {
		res, err := repository.ExemptedService.DBClient.ExemptUserVanity(user.ID, "staff")
		if err != nil || !res {
			msg := "Failed to add vanity exemption"
			if _, e := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg}); e != nil {
				fmt.Println("Failed to edit deferred interaction:", e)
			}
			return
		}
	} else {
		res, err := repository.ExemptedService.DBClient.ExemptUserVanity(user.ID, "")
		if err != nil || !res {
			msg := "Failed to add vanity exemption"
			if _, e := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg}); e != nil {
				fmt.Println("Failed to edit deferred interaction:", e)
			}
			return
		}
	}

	newRoleIds := strings.Split(os.Getenv("SUPPORTER_ROLE_IDS"), ",")
	if len(newRoleIds) == 0 {
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

	updatedUserRoles := append(member.Roles, newRoleIds...)
	if !slices.Equal(member.Roles, updatedUserRoles) {
		_, err := s.GuildMemberEdit(i.GuildID, user.ID, &discordgo.GuildMemberParams{
			Roles: &updatedUserRoles,
		})
		if err != nil {
			msg := "Failed to update user roles."
			if _, e := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg}); e != nil {
				fmt.Println("Failed to edit deferred interaction:", e)
			}
			return
		}
		successMsg := fmt.Sprintf("Vanity added for user <@%s>. They now have the supporter role.", user.ID)
		if _, e := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &successMsg}); e != nil {
			fmt.Println("Failed to edit deferred interaction:", e)
		}
		return
	}
}