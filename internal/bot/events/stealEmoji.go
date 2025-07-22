package events

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Sush1sui/fns-go/internal/common"
	"github.com/bwmarrin/discordgo"
)

func OnStealEmoji(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.ID == "" || !strings.HasPrefix(m.Content, "!se") {
		return
	}

	re := regexp.MustCompile(`<a?:(\w+):(\d+)>`)
	emoji := re.FindStringSubmatch(m.Content)
	if emoji == nil {
		fmt.Println("No emoji found in message")
		return
	}

	isAnimated := strings.HasPrefix(emoji[0], "<a:")
	emojiName := emoji[1]
	emojiID := emoji[2]
	emojiURL := fmt.Sprintf("https://cdn.discordapp.com/emojis/%s.%s", emojiID, map[bool]string{true: "gif", false: "png"}[isAnimated])
	fmt.Println(emoji)

	assetsDir := filepath.Join(".", "assets")
	if _, err := os.Stat(assetsDir); os.IsNotExist(err) {
		os.MkdirAll(assetsDir, 0755)
	}

	fileName := filepath.Base(emojiURL)
	destination := filepath.Join(assetsDir, fileName)

	// Download the emoji file
	if err := common.DownloadFile(emojiURL, destination); err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed to download emoji: "+err.Error())
		return
	}

	// Read the downloaded file
	fileBuffer,  err := os.ReadFile(destination)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed to read downloaded emoji file: "+err.Error())
		return
	}

	// Prepare image as base64 data URI
	mimeType := "image/png"
	if isAnimated {
			mimeType = "image/gif"
	}
	base64Image := base64.StdEncoding.EncodeToString(fileBuffer)
	dataURI := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Image)

	// Upload the emoji to the server
	newEmoji, err := s.GuildEmojiCreate(m.GuildID, &discordgo.EmojiParams{
		Name:  emojiName,
		Image: dataURI,
		Roles: []string{}, // You can specify roles if needed
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Failed to upload emoji: %v", err))
		os.Remove(destination)
		return
	}

	// Respond to the user after the emoji is uploaded
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("**<:pepehorny:1303965571528003606> Successfully stolen: <:%s:%s>**", newEmoji.Name, newEmoji.ID))

	// Delete the file from the assets folder after the upload
	os.Remove(destination)
}