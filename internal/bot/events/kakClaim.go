package events

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/Sush1sui/fns-go/internal/common"
	"github.com/bwmarrin/discordgo"
)

var kakClaimCommands = []string{"$kak claim", "$kc", "$tc", "$trash claim"}


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

	common.KakClaimMu.Lock()
	// Stop and clean up previous timers if they exist
	if timer, exists := common.KakClaimTimeoutMap[m.Author.ID]; exists {
		timer.Stop()
		delete(common.KakClaimTimeoutMap, m.Author.ID)
	}
	if ticker, exists := common.KakClaimIntervalMap[m.Author.ID]; exists {
		ticker.Stop()
		delete(common.KakClaimIntervalMap, m.Author.ID)
	}
	common.KakClaimMu.Unlock()

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
		common.KakClaimMu.Lock()
		delete(common.KakClaimTimeoutMap, m.Author.ID)
		delete(common.KakClaimIntervalMap, m.Author.ID)
		common.KakClaimMu.Unlock()
	})

	common.KakClaimMu.Lock()
	common.KakClaimTimeoutMap[m.Author.ID] = timeout
	common.KakClaimIntervalMap[m.Author.ID] = ticker
	common.KakClaimMu.Unlock()

	go func() {
		for i := remainingTime - 1; i > 0; i-- {
			<-ticker.C
			// Update the countdown message
			_, err := s.ChannelMessageEdit(m.ChannelID, replyMsg.ID, fmt.Sprintf("**<@%s> can kak/trash claim in %d seconds!**", m.Author.ID, i))
			if err != nil {
				fmt.Println("Error editing kak claim message:", err)
				return
			}
		}
	}()
}