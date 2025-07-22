package events

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func OnAutoReact(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if os.Getenv("CHANNEL_ID_SELFIE") == m.ChannelID && len(m.Attachments) > 0 {
		go func() {
			err := s.MessageReactionAdd(m.ChannelID, m.ID, "Check_White_FNS:1310274014102687854")
			if err != nil {
				fmt.Println("Failed to add reaction: " + err.Error())
			}
		}()
		go func() {
			err := s.MessageReactionAdd(m.ChannelID, m.ID, "pixelheart:1310424521421099113")
			if err != nil {
				fmt.Println("Failed to add reaction: " + err.Error())
			}
		}()
	}

	
	if strings.Contains(strings.ReplaceAll(strings.ToLower(m.Content), " ", ""), "hahaha") {
		go func() {
			err := s.MessageReactionAdd(m.ChannelID, m.ID, "😂")
			if err != nil {
				fmt.Println("Failed to add reaction: " + err.Error())
			}
		}()
	}

	if strings.Contains(strings.ReplaceAll(strings.ToLower(m.Content), " ", ""), "hahaha") {
		go func() {
			err := s.MessageReactionAdd(m.ChannelID, m.ID, "SushiRolling:1293411594621157458")
			if err != nil {
				fmt.Println("Failed to add reaction: " + err.Error())
			}
		}()
	}
}