package events

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/Sush1sui/fns-go/internal/common"
	"github.com/bwmarrin/discordgo"
)

var (
	kakClaimCommands   = []string{"$kak claim", "$kc", "$tc", "$trash claim"}
	kakClaimTimeoutMap  = make(map[string]*time.Timer)
	kakClaimIntervalMap = make(map[string]*time.Ticker)
	kakClaimMu          = sync.Mutex{}
)

func OnKakClaim(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}
	if strings.ToLower(m.ChannelID) != os.Getenv("MUDAE_CHANNEL_ID") {
		return
	}

	content := strings.ToLower(m.Content)
	found := slices.Contains(kakClaimCommands, content)
	if !found {
		return
	}

	kakClaimMu.Lock()
	// Stop and clean up previous timers if they exist
	if timer, exists := kakClaimTimeoutMap[m.Author.ID]; exists {
		timer.Stop()
		delete(kakClaimTimeoutMap, m.Author.ID)
	}
	if ticker, exists := kakClaimIntervalMap[m.Author.ID]; exists {
		ticker.Stop()
		delete(kakClaimIntervalMap, m.Author.ID)
	}
	kakClaimMu.Unlock()

	// Start a new timer for the user
	remainingTime := int(common.KakClaimTimer.Seconds())

	replyMsg, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("**<@%s> can kak/trash claim in %d seconds!**", m.Author.ID, remainingTime))
	if err != nil {
		fmt.Println("Error sending kak claim message:", err)
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	timeout := time.AfterFunc(common.KakClaimTimer, func() {
		ticker.Stop()
		// Delete the countdown message after the timer expires
		_ = s.ChannelMessageDelete(m.ChannelID, replyMsg.ID)
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("**<@%s> you can now kak/trash claim!**", m.Author.ID))
		if err != nil {
			fmt.Println("Error sending kak claim ready message:", err)
		}
		kakClaimMu.Lock()
		delete(kakClaimTimeoutMap, m.Author.ID)
		delete(kakClaimIntervalMap, m.Author.ID)
		kakClaimMu.Unlock()
	})

	kakClaimMu.Lock()
	kakClaimTimeoutMap[m.Author.ID] = timeout
	kakClaimIntervalMap[m.Author.ID] = ticker
	kakClaimMu.Unlock()

	go func() {
		for i := remainingTime - 1; i > 0; i-- {
			_, ok := <-ticker.C
			if !ok {
				return // Channel closed, exit the goroutine
			}
			// Update the countdown message
			_, err := s.ChannelMessageEdit(m.ChannelID, replyMsg.ID, fmt.Sprintf("**<@%s> can kak/trash claim in %d seconds!**", m.Author.ID, i))
			if err != nil {
				fmt.Println("Error editing kak claim message:", err)
				return
			}
		}
	}()
}