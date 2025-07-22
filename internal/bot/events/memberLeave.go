package events

import (
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func OnMemberLeave(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	if m.User.ID == s.State.User.ID || m.User.Bot {
		return
	}

	welcomeChannel, err := s.Channel(os.Getenv("WELCOME_CHANNEL_ID"))
	if err != nil || welcomeChannel == nil {
		return
	}

	var lastId string
	fetchCount := 0

	for fetchCount < 10 {
		messages, err := s.ChannelMessages(welcomeChannel.ID, 100, lastId, "", "")
		if err != nil || len(messages) == 0 {
			break
		}


		for _, msg := range messages {
			if msg.Author != nil && msg.Author.ID == s.State.User.ID &&
				len(msg.Embeds) > 0 &&
				strings.Contains(msg.Embeds[0].Title, "<@"+m.User.ID+">") {
					_ = s.ChannelMessageDelete(welcomeChannel.ID, msg.ID)
					break
				}
		}

		lastId = messages[len(messages)-1].ID
		fetchCount++
	}
}