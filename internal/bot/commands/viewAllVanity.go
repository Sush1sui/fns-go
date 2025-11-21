package commands

import (
	"fmt"

	"github.com/Sush1sui/fns-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func VanityViewAll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil { return }

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Flags: discordgo.MessageFlagsEphemeral,
        },
    })
    
	allExemptedUsers, err := repository.ExemptedService.DBClient.GetAllExemptedUsers()
	if err != nil {
		fmt.Println("Something went wrong with getting all exempted users")
		msg := "❌ | An error occurred while fetching exempted users."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
		return
	}

	mentionString := ""
	for idx, exemptedUser := range allExemptedUsers {
		suffix := ""
		if exemptedUser.Expiration.IsZero() {
			suffix = " — Permanent"
		} else {
			// show expiration time in UTC for clarity
			// print in unix timestamp format for discord to render nicely
			suffix = fmt.Sprintf(" — Expires: <t:%d:R>", exemptedUser.Expiration.Unix())
		}

		if idx == len(allExemptedUsers)-1 {
			mentionString += fmt.Sprintf("<@%s>%s", exemptedUser.UserID, suffix)
		} else {
			mentionString += fmt.Sprintf("<@%s>%s\n", exemptedUser.UserID, suffix)
		}
	}

	msg := fmt.Sprintf("📋 | **Exempted Users (%d):**\n%s", len(allExemptedUsers), mentionString)
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
}