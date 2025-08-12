package events

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func OnMemberJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	if m.User.ID == s.State.User.ID || m.User.Bot { return }

	if m.Member.User.ID != "1258348384671109120" {
		s.ChannelMessageSendEmbed(os.Getenv("WELCOME_CHANNEL_ID"), &discordgo.MessageEmbed{
			Color: 0xffffff,
			Title: "-ˏˋ⋆ ᴡ ᴇ ʟ ᴄ ᴏ ᴍ ᴇ ⋆ˊˎ-",
			Description: fmt.Sprintf(
					"Hello <@%s>! Welcome to **Finesse**.\n\n"+
							"Please make sure you head to <#1303919197629321308> before chatting.\n"+
							"On top of that, please go to <#1292714443351785502> to set up your profile.\n\n"+
							"└─── we hope you enjoy your stay in here!──➤",
					m.User.ID,
			),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://cdn.discordapp.com/attachments/1293239740404994109/1310247866970800209/tl.png?ex=674486ea&is=6743356a&hm=e89c4eb171c56724f5e7b1702d85acfa208a03b47c65632de850af79fe826a8c&",
			},
		})
	}

  	lastId := ""
	fetchCount := 0

	for fetchCount < 10 {
		// For those saying "Oh you leaked your channel id"
		// everyone on discord can get any channel ids and messages
		messages, err := s.ChannelMessages("1292442961514201150", 100, lastId, "", "")
		if err != nil || len(messages) == 0 { break }

		for _, msg := range messages {
			if msg.Author != nil && msg.Author.ID == "1292751642822967319" && len(msg.Embeds) > 0 && strings.Contains(msg.Embeds[0].Description, "1258348384671109120") {
				_ = s.ChannelMessageDelete("1292442961514201150", msg.ID)
				break
			}
		}

		lastId = messages[len(messages)-1].ID
		fetchCount++
	}

	if m.User.ID == "1258348384671109120" {
		err := s.GuildMemberRoleAdd(m.GuildID, m.Member.User.ID, "1321872792089526372")
		if err != nil {
			member, err := s.GuildMember(m.GuildID, "982491279369830460")
			if err != nil { return }

			dmChannel, err := s.UserChannelCreate(member.User.ID)
			if err != nil { return }

			s.ChannelMessageSend(dmChannel.ID, "Failed to add mudae role to Dane.")
		}
	}
}