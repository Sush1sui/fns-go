package events

import (
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func OnMemberLeave(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	if m.User == nil || s.State == nil || s.State.User == nil {
		return
	}
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
		if err != nil || len(messages) == 0 { break }
		if len(messages) < 1 { break }

		for _, msg := range messages {
			if msg.Author != nil && msg.Author.ID == s.State.User.ID &&
				len(msg.Embeds) > 0 &&
				strings.Contains(msg.Embeds[0].Description, "Welcome") {
					_ = s.ChannelMessageDelete(welcomeChannel.ID, msg.ID)
					break
				}
		}

		lastId = messages[len(messages)-1].ID
		fetchCount++
	}

	lastId = ""
	fetchCount = 0

	for fetchCount < 10 {
		// For those saying "Oh you leaked your channel id"
		// everyone on discord can get any channel ids and messages
		messages, err := s.ChannelMessages("1292442961514201150", 100, lastId, "", "")
		if err != nil || len(messages) == 0 { break }
		if len(messages) < 1 { break }

		for _, msg := range messages {
			if msg.Author != nil && msg.Author.ID == "1292751642822967319" && len(msg.Embeds) > 0 && strings.Contains(msg.Embeds[0].Description, "1258348384671109120") {
				_ = s.ChannelMessageDelete("1292442961514201150", msg.ID)
				break
			}
		}

		lastId = messages[len(messages)-1].ID
		fetchCount++
	}
}