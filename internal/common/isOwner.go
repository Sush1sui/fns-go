package common

import "github.com/bwmarrin/discordgo"

// Checks if the user is the server owner
func IsGuildOwner(guild *discordgo.Guild, userID string) bool {
	return guild.OwnerID == userID
}