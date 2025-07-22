package events

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func OnMemberJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	if m.User.ID == s.State.User.ID || m.User.Bot {
		return
	}

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