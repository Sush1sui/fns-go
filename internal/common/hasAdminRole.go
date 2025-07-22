package common

import "github.com/bwmarrin/discordgo"

// Checks if the user has a role with admin permissions
func HasAdminRole(guild *discordgo.Guild, member *discordgo.Member) bool {
	for _, userRoleID := range member.Roles {
		for _, guildRole := range guild.Roles {
			if guildRole.ID == userRoleID && (guildRole.Permissions&discordgo.PermissionAdministrator) != 0 {
				return true
			}
		}
	}
	return false
}