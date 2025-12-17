package helpers

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Sush1sui/fns-go/internal/common"
	"github.com/Sush1sui/fns-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

var rainbowTransition = []string{
	"#FF0000", // 0° Red
	"#FF7F00", // 15° Orange
	"#FFBF00", // 30° Golden yellow
	"#FFFF00", // 45° Yellow
	"#BFFF00", // 60° Yellow-green
	"#7FFF00", // 75° Lime
	"#3FFF00", // 90° Spring green
	"#00FF00", // 105° Green
	"#00FF3F", // 120° Green-cyan
	"#00FF7F", // 135° Aqua-green
	"#00FFBF", // 150° Turquoise
	"#00FFFF", // 165° Cyan
	"#00BFFF", // 180° Sky blue
	"#007FFF", // 195° Azure
	"#003FFF", // 210° Blue
	"#0000FF", // 225° Deep blue
	"#3F00FF", // 240° Indigo
	"#7F00FF", // 255° Violet
	"#BF00FF", // 270° Purple
	"#FF00FF", // 285° Magenta
	"#FF00BF", // 300° Hot pink
	"#FF007F", // 315° Rose
	"#FF003F", // 330° Scarlet
	"#FF0000", // 345° Back to Red
}

func ScanForVanityLinks(s *discordgo.Session) {
	guild, err := s.State.Guild(os.Getenv("GUILD_ID")) // Replace with your guild ID
	if err != nil {
		fmt.Println("Error fetching guild for vanity:", err)
		return
	}

	supporterLink := os.Getenv("SUPPORTER_LINK")
	roleIDs := strings.Split(os.Getenv("SUPPORTER_ROLE_IDS"), ",")
	if len(roleIDs) < 2 {
		fmt.Println("SUPPORTER_ROLE_IDS environment variable must contain at least 2 role IDs")
		return
	}
	supporterRoleID := roleIDs[0]
	exemptedRoleID := roleIDs[1]

	if supporterLink == "" || supporterRoleID == "" {
		fmt.Println("SUPPORTER_LINK or SUPPORTER_ROLE_IDS environment variable is not set")
		return
	}

	fmt.Println("Scanning for vanity links")
	// Fetch all members in the guild
	var allMembers []*discordgo.Member
	var after string
	for {
		members, err := s.GuildMembers(guild.ID, after, 1000)
		if err != nil {
			fmt.Println("Error fetching members:", err)
			break
		}
		if len(members) == 0 {
			break
		}
		allMembers = append(allMembers, members...)
		after = members[len(members)-1].User.ID
		if len(members) < 1000 {
			break
		}
	}

	exemptedUsers, err := repository.ExemptedService.DBClient.GetAllExemptedUsers()
	if err != nil {
		fmt.Println("Error fetching exempted users:", err)
		return
	}

	// filter out expired vanity users locally — MongoDB TTL will remove expired docs automatically
	now := time.Now().UTC()
	filtered := exemptedUsers[:0] // reuse backing array
	for _, u := range exemptedUsers {
		if u.UserID == "1258348384671109120" {continue}

		if u.Expiration.IsZero() || u.Expiration.After(now) {
			filtered = append(filtered, u)
		} else {
			// detected expired entry; TTL index (server-side) will delete this document shortly
			fmt.Println("(Expired vanity) Detected expired user (TTL will remove):", u.UserID)
		}
	}
	exemptedUsers = filtered

	// Create a map for O(1) lookups
	exemptedUserIDs := make(map[string]bool)
	for _, user := range exemptedUsers {
		exemptedUserIDs[user.UserID] = true
	}

	var supporterRole *discordgo.Role
	for _, role := range guild.Roles {
		if role.ID == supporterRoleID {
			supporterRole = role
			break
		}
	}
	if supporterRole == nil {
		fmt.Println("Supporter role not found")
		return
	}

	colorIndex := 0
	currentColor := rainbowTransition[0]

	var supporterChannel *discordgo.Channel
	supporterChannelID := os.Getenv("SUPPORTER_CHANNEL_ID")
	for _, channel := range guild.Channels {
		if channel.ID == supporterChannelID {
			supporterChannel = channel
			break
		}
	}
	if supporterChannel == nil {
		fmt.Println("Supporter channel not found")
		return
	}

	// This helps to make O(1) presence lookups
	presenceMap := make(map[string]*discordgo.Presence)
	for _, p := range guild.Presences {
		if p.User != nil {
			presenceMap[p.User.ID] = p
		}
	}

	for _, member := range allMembers {
		// skip bots
		if member.User.Bot { continue }

		// Find the presence for this member
		presence := presenceMap[member.User.ID]
		if presence == nil {
		continue // No presence info for this member
		}

		// skip staff members
		isStaff := false
		for _, rid := range member.Roles {
			if slices.Contains(common.StaffRoleIDs, rid) {
				isStaff = true
				break
			}
		}
		if isStaff { continue }

		// Find the custom status activity with the desired state
		var customStatus string
		for _, activity := range presence.Activities {
			if activity.State == supporterLink {
				customStatus = activity.State
				break
			}
		}

		includesSupporterLink := strings.Contains(customStatus, supporterLink) || exemptedUserIDs[member.User.ID]
		hasSupporterRole := slices.Contains(member.Roles, supporterRole.ID)


		if includesSupporterLink && hasSupporterRole {
			fmt.Println("Member already has the role and status:", member.User.Username)
			continue
		}

		// add or remove the role based on the link
		if(includesSupporterLink && !hasSupporterRole) {
			member.Roles = append(member.Roles, supporterRole.ID)
			_, err := s.GuildMemberEdit(guild.ID, member.User.ID, &discordgo.GuildMemberParams{
				Roles: &member.Roles,
			})
			if err != nil {
				fmt.Println("Error adding supporter role to member:", err)
			}

			_, err = s.ChannelMessageSendEmbed(supporterChannel.ID, &discordgo.MessageEmbed{
					Title: "Thank you for supporting **Finesse!**",
					Description: fmt.Sprintf(
							"<@%s> updated their status with our vanity link `discord.gg/finesseph` and earned the %s role!\n\n"+
									"> Perks:\n"+
									"- Image & Embed Link Perms\n"+
									"- 1.5x Level Boost\n"+
									"- Color Name <#1310451488975224973>\n",
							member.User.ID, supporterRole.Name,
					),
					Image: &discordgo.MessageEmbedImage{
							URL: "https://cdn.discordapp.com/attachments/1293239740404994109/1310449852349681704/image.png",
					},
					Color: func() int {
							c, _ := strconv.ParseInt(strings.Replace(currentColor, "#", "", 1), 16, 32)
							return int(c)
					}(),
					Footer: &discordgo.MessageEmbedFooter{
							Text: "*Note: Perks will be revoked if you remove the status.*",
					},
			})
			if err != nil {
				fmt.Println("Error sending supporter message:", err)
			}

			currentColor = rainbowTransition[(colorIndex + 1) % len(rainbowTransition)]
		} else if !includesSupporterLink && hasSupporterRole {
			// remove the role
			rolesToRemove := map[string]struct{}{
				supporterRoleID: {},
				exemptedRoleID: {},
			}
			newRoles := make([]string, 0, len(member.Roles))
			for _, roleId := range member.Roles {
				if _, remove := rolesToRemove[roleId]; !remove {
					newRoles = append(newRoles, roleId)
				}
			}
			member.Roles = newRoles
			_, err := s.GuildMemberEdit(guild.ID, member.User.ID, &discordgo.GuildMemberParams{
				Roles: &member.Roles,
			})
			if err != nil {
				fmt.Println("Error removing supporter role from member:", err)
			}
		}
	}

	fmt.Println("Vanity link scan completed")
}