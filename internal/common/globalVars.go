package common

import (
	"os"
	"strings"
	"time"
)

var StickyChannels = map[string]struct{}{}
var GuildID = ""
var StaffRoleIDs = []string{}
var ChannelExceptionIDs = []string{}
var PrivilegedRoleIDs = []string{}

func InitializeGlobalVars() {
	StaffRoleIDs = append(StaffRoleIDs, strings.Split(os.Getenv("STAFF_ROLE_IDS"), ",")...)
	ChannelExceptionIDs = append(ChannelExceptionIDs, strings.Split(os.Getenv("CHANNEL_EXCEPTIONS_IDS"), ",")...)
	PrivilegedRoleIDs = append(PrivilegedRoleIDs, strings.Split(os.Getenv("PRIVILEDGED_ROLE_IDS"), ",")...)
	GuildID = os.Getenv("GUILD_ID")
}

var Keywords = []string{"http://", "https://", "www.", "discord.gg/"}

var KakClaimTimer = 15 * time.Second
