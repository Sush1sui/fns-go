package events

import (
	"fmt"
	"strings"

	"github.com/Sush1sui/fns-go/internal/common"
	"github.com/Sush1sui/fns-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func OnNicknameRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author == nil {
		return
	}
	if m.Author.Bot{
		return
	}
	if m.ChannelID != "1310583941287379116" {
		return
	}
	if !strings.HasPrefix(strings.ToLower(m.Content), "!rn") {
		return
	}
	fmt.Println("Received nickname request from user:", m.Author.ID)

	nicknameRequest := strings.TrimSpace(m.Content[len("!rn"):])
	if nicknameRequest == "" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Please provide a nickname request after the command.")
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
		return
	}

	if len(common.StaffRoleIDs) == 0 {
		fmt.Println("No staff roles configured")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Nickname Request",
		Description: fmt.Sprintf("**<@%s> has requested the nickname: %s**.\n\n<@&%s>, please approve or decline.", m.Author.ID, nicknameRequest, common.StaffRoleIDs[0]),
		Color:       0xFFFFFF, // White color
	}

	approvalChannel, err := s.State.Channel("1310273100583276544")
	if err != nil {
		fmt.Println("Error fetching approval channel:", err)
		return
	}

	content := fmt.Sprintf("<@&%s>", common.StaffRoleIDs[0])
	message, err := s.ChannelMessageSendComplex(approvalChannel.ID, &discordgo.MessageSend{
			Content: content, // This will mention the staff role
			Embed:   embed,
	})
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	err = s.MessageReactionAdd(message.ChannelID, message.ID, "Check_White_FNS:1310274014102687854")
	if err != nil {
		fmt.Println("Error adding approve reaction:", err)
		return
	}
	err = s.MessageReactionAdd(message.ChannelID, message.ID, "No:1310633209519669290")
	if err != nil {
		fmt.Println("Error adding deny reaction:", err)
		return
	}

	// Create nickname request in the database
	_, err = repository.NicknameRequestService.DBClient.CreateNicknameRequest(nicknameRequest, m.Author.ID, m.ID, m.ChannelID, approvalChannel.ID, message.ID)
	if err != nil {
		fmt.Println("Error creating nickname request:", err)
		return
	}

	err = repository.NicknameRequestService.DBClient.SetupNicknameRequestCollector(s, message, nicknameRequest)
	if err != nil {
		fmt.Println("Error setting up nickname request collector:", err)
		return
	}
	fmt.Printf("Nickname request created for user %s: %s\n", m.Author.ID, nicknameRequest)
}