package common

import (
	"os"
	"strings"
	"sync"
	"time"
)

var StickyChannels = map[string]struct{}{}
var StaffRoleIDs = strings.Split(os.Getenv("STAFF_ROLE_IDS"), ",")
var ChannelExceptionIDs = strings.Split(os.Getenv("CHANNEL_EXCEPTIONS_IDS"), ",")
var PrivilegedRoleIDs = strings.Split(os.Getenv("PRIVILEDGED_ROLE_IDS"), ",")

var Keywords = []string{"http", "www.", "discord.gg/"}

var KakClaimTimer = 15 * time.Second
var KakClaimTimeoutMap = make(map[string]*time.Timer)
var KakClaimIntervalMap = make(map[string]*time.Ticker)
var KakClaimMu = sync.Mutex{}