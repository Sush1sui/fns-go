package common

import (
	"os"
	"strings"
)

var StickyChannels = map[string]struct{}{}
var StaffRoleIDs = strings.Split(os.Getenv("STAFF_ROLE_IDS"), ",")
var ChannelExceptionIDs = strings.Split(os.Getenv("CHANNEL_EXCEPTIONS_IDS"), ",")
var PrivilegedRoleIDs = strings.Split(os.Getenv("PRIVILEDGED_ROLE_IDS"), ",")

var Keywords = []string{"http", "www.", "discord.gg/"}



