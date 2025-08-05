package events

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var charactersSet = make(map[string]struct{})

func init() {
    file, err := os.Open("internal/common/characters.json")
    if err != nil {
        fmt.Println("Error opening characters.json:", err)
        return
    }
    defer file.Close()

    var characters []string
    if err := json.NewDecoder(file).Decode(&characters); err != nil {
        fmt.Println("Error decoding characters.json:", err)
        return
    }
    for _, c := range characters {
        charactersSet[strings.ToLower(c)] = struct{}{}
    }
    fmt.Println("Characters loaded successfully from characters.json with", len(charactersSet), "entries.")
}

func OnSnipeMudae(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != "432610292342587392" {
		return
	}
	if m.ChannelID != os.Getenv("MUDAE_CHANNEL_ID") {
		return
	}

	if len(m.Embeds) == 0 || m.Embeds[0] == nil {
		return
	}
	embed := m.Embeds[0]

	if embed == nil || embed.Footer == nil || embed.Author == nil {
		return
	}

	if !strings.Contains(strings.ToLower(embed.Footer.Text), "belongs to") {
		if _, ok := charactersSet[strings.ToLower(embed.Author.Name)]; ok {
			fmt.Println("Top character found:", embed.Author.Name)
			vipUsers := strings.Split(os.Getenv("SNIPER_VIP_USERS"), ",")
			for _, id := range vipUsers {
				go func(userID string) {
					user, err := s.User(userID)
					if err != nil || user == nil {
						fmt.Println("Error fetching user:", err)
						return
					}
					dmChannel, err := s.UserChannelCreate(userID)
					if err != nil {
						fmt.Println("Error creating DM channel:", err)
						return
					}

					// Construct the jump link
					messageURL := fmt.Sprintf("https://discord.com/channels/%s/%s/%s", m.GuildID, m.ChannelID, m.ID)
					content := fmt.Sprintf("# A top character `%s` has appeared! Click here to jump to the message: %s", embed.Author.Name, messageURL)

					_, err = s.ChannelMessageSendComplex(dmChannel.ID, &discordgo.MessageSend{
						Content: content,
						Embed:   embed,
					})
					if err != nil {
						fmt.Println("Error sending DM:", err)
						return
					}
				}(id)
			}
		}
	}
}