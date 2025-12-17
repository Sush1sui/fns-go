package events

import (
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

var memberBoostCache = make(map[string]int64) // userID -> PremiumSince.Unix()

func OnMemberBoost(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	if m.User == nil || s.State == nil || s.State.User == nil {
		return
	}
	if m.User.ID == s.State.User.ID || m.User.Bot { return }

	boostChannel, err := s.Channel(os.Getenv("BOOST_CHANNEL_ID"))
	if err != nil || boostChannel == nil {
		fmt.Println("Error fetching boost channel:", err)
		return
	}

	serverLogChannel, err := s.Channel(os.Getenv("SERVER_LOG_CHANNEL_ID"))
	if err != nil || serverLogChannel == nil {
		fmt.Println("Error fetching server log channel:", err)
		return
	}

	// Detect boost: compare previous and current PremiumSince
	prevBoost := memberBoostCache[m.User.ID]
	newBoost := int64(0)
	if m.PremiumSince != nil {
		newBoost = m.PremiumSince.Unix()
	}

	if newBoost > 0 && (prevBoost == 0 || prevBoost < newBoost) {
		embed := &discordgo.MessageEmbed{
			Title:       "Thank you for the server boost!",
			Color:       0xff73fa,
			Description: "** We truly appreciate your support and all you do to help make this community even better! Sending you all our love and gratitude!**\n\n" +
					"> **Perks**\n" +
					"- Receive <@&1292420325002448930> role\n" +
					"- Custom Onigiri Color Role <#1303919788342382615>\n" +
					"- Nickname perms\n" +
					"- Soundboard\n" +
					"- Image and Embed Links perms\n" +
					"- External Emoji & Sticker\n" +
					"- 2.0x Level Boost",
			Image: &discordgo.MessageEmbedImage{
				URL: "https://cdn.discordapp.com/attachments/1303917209101406230/1310235247341867028/tyfdboost.gif?ex=67447b29&is=674329a9&hm=9a9a997a62dc63ea9db4071304ab3489d56149113fa7d5f0479d42aea1c3f1ed&",
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}

		_, err := s.ChannelMessageSendComplex(boostChannel.ID, &discordgo.MessageSend{
			Embed: embed,
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
			},
			Content: fmt.Sprintf("# <@%s> HAS BOOSTED THE SERVER", m.User.ID),
		})
		if err != nil {
			fmt.Println("Error sending boost message:", err)
			return
		}

		// Log the boost in the server log channel
		_, err = s.ChannelMessageSend(serverLogChannel.ID, fmt.Sprintf("**%s** has boosted the server! Thank you for your support!", m.User.Username))
		if err != nil {
			fmt.Println("Error logging boost message:", err)
			return
		}

		// Update the cache
		memberBoostCache[m.User.ID] = newBoost
	}
}

func SyncMemberBoostCache(s *discordgo.Session, guildID string) {
    fmt.Println("Syncing member boost cache...")
    after := ""
    for {
        members, err := s.GuildMembers(guildID, after, 1000)
        if err != nil {
            fmt.Println("Error fetching guild members for boost cache:", err)
            break
        }
        if len(members) == 0 {
            break
        }
        for _, m := range members {
            if m.PremiumSince != nil {
                memberBoostCache[m.User.ID] = m.PremiumSince.Unix()
            }
            after = m.User.ID
        }
        if len(members) < 1000 {
            break
        }
    }
    fmt.Printf("Synced %d boosted members to cache.\n", len(memberBoostCache))
}