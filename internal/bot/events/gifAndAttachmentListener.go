package events

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/Sush1sui/fns-go/internal/common"
	"github.com/bwmarrin/discordgo"
)

var (
    batchDeleteMap   = make(map[string][]string) // channelID -> messageIDs
    batchDeleteTimes = make(map[string]time.Time)
    batchDeleteMu    sync.Mutex
    batchThreshold   = 10
    batchWindow      = 5 * time.Second
)

func OnGifAndAttachment(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot || slices.Contains(common.ChannelExceptionIDs, m.ChannelID) {
		return
	}
	if m.Member == nil {
		return
	}

	// Check if the user has the necessary permissions
	if common.IsGuildOwner(&discordgo.Guild{ID: m.GuildID}, m.Author.ID) {
		return
	}
	if !common.HasAdminRole(&discordgo.Guild{ID: m.GuildID}, m.Member) {
		return
	}
	for _, id := range append(common.PrivilegedRoleIDs, common.StaffRoleIDs...) {
		if len(m.Member.Roles) > 0 && slices.Contains(m.Member.Roles, id) {
			return
		}
	}

	hasAttachments := len(m.Attachments) > 0
	
	includesLinks := false
	for _, keyword := range common.Keywords {
		if strings.Contains(strings.ToLower(m.Content), keyword) {
			includesLinks = true
			break
		}
	}

	if hasAttachments || includesLinks {
		batchDeleteMu.Lock()
		batchDeleteMap[m.ChannelID] = append(batchDeleteMap[m.ChannelID], m.ID)
		now := time.Now()
		if _, ok := batchDeleteTimes[m.ChannelID]; !ok {
			batchDeleteTimes[m.ChannelID] = now
		}
		// If threshold reached and within window, batch delete
		if len(batchDeleteMap[m.ChannelID]) >= batchThreshold && now.Sub(batchDeleteTimes[m.ChannelID]) <= batchWindow {
			ids := batchDeleteMap[m.ChannelID]
			batchDeleteMap[m.ChannelID] = nil
			batchDeleteTimes[m.ChannelID] = now
			batchDeleteMu.Unlock()
			go func() {
				for _, id := range ids {
					err := s.ChannelMessageDelete(m.ChannelID, id)
					if err != nil {
						fmt.Println("Failed to batch delete message:", err)
					}
				}
			}()
		} else {
			batchDeleteMu.Unlock()
			// If not enough for batch, delete single message as usual
			go func() {
				err := s.ChannelMessageDelete(m.ChannelID, m.ID)
				if err != nil {
					fmt.Println("Failed to delete message with attachments or links:", err)
				}
			}()
		}
	}
}
