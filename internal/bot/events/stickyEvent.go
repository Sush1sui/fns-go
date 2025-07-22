package events

import (
	"fmt"
	"sync"

	"github.com/Sush1sui/fns-go/internal/common"
	"github.com/Sush1sui/fns-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

var stickyLocks = make(map[string]*sync.Mutex)
var stickyLocksMu sync.Mutex

func getStickyLock(channelId string) *sync.Mutex {
	stickyLocksMu.Lock()
	defer stickyLocksMu.Unlock()
	if stickyLocks[channelId] == nil {
		stickyLocks[channelId] = &sync.Mutex{}
	}
	return stickyLocks[channelId]
}

func OnSticky(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return // Ignore messages from the bot itself
	}

	lock := getStickyLock(m.ChannelID)
	lock.Lock()
	defer lock.Unlock()
	
	// Check if the channel is a sticky channel
	if _, ok := common.StickyChannels[m.ChannelID]; !ok {
			return
	}

	stickyChannel := repository.StickyService.DBClient.GetStickyChannel(m.ChannelID)
	if stickyChannel == nil {
		fmt.Println("Sticky channel not found in database.")
		return
	}

	// If there is a last sticky message, delete it
	if stickyChannel.LastStickyMessageId != "" {
		err := s.ChannelMessageDelete(m.ChannelID, stickyChannel.LastStickyMessageId)
		if err != nil {
			fmt.Println("Error deleting last sticky message:", err)
		}
		fmt.Println("Deleted last sticky message:", stickyChannel.LastStickyMessageId)
	}

	// fetch last 2 messages sent in channel (if there are any)
	messages, err := s.ChannelMessages(m.ChannelID, 2, "", "", "")
	if err != nil {
		fmt.Println("Error retrieving messages from channel:", err)
		return
	}
	recentMessageId := messages[1].ID

	embed := &discordgo.MessageEmbed{
		Title:       "Stickied Message",
		Description: stickyChannel.StickyMessage,
		Color:       0xFFFFFF, //white
	}

	newStickyMessge, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		fmt.Println("Error sending sticky message:", err)
		return
	}

	// Save the new sticky message ID
	_, err = repository.StickyService.DBClient.UpdateStickyMessageId(m.ChannelID, newStickyMessge.ID, recentMessageId)
	if err != nil {
		fmt.Println("Error updating sticky message ID:", err)
		return
	}
}