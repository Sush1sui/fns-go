package mongodb

import (
	"context"
	"fmt"
	"os"

	"github.com/Sush1sui/fns-go/internal/model"
	"github.com/Sush1sui/fns-go/internal/repository"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var approveEmoji = os.Getenv("APPROVE_EMOJI")
var denyEmoji = os.Getenv("DENY_EMOJI")
var approvedEmojiID = os.Getenv("APPROVED_EMOJI_ID")
var deniedEmojiID = os.Getenv("DENIED_EMOJI_ID")

func (c *MongoClient) GetAllNicknameRequests() ([]*model.NicknameRequest, error) {
	var requests []*model.NicknameRequest

	cursor, err := c.Client.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error finding nickname requests: %v", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var request model.NicknameRequest
		if err := cursor.Decode(&request); err != nil {
			return nil, fmt.Errorf("error decoding nickname request: %v", err)
		}
		requests = append(requests, &request)
	}
	return requests, nil
}

func (c *MongoClient) CreateNicknameRequest(nickname, userId, userMessageId, userChannelId, staffChannelId, staffMessageId string) (*model.NicknameRequest, error) {
	request := &model.NicknameRequest{
		Nickname:        	nickname,
		UserID:         	userId,
		UserMessageID:  	userMessageId,
		UserChannelID:  	userChannelId,
		StaffChannelID: 	staffChannelId,
		StaffMessageID: 	staffMessageId,
		Reactions: []struct {
			Emoji string "bson:\"emoji\""
		}{
			{Emoji: approveEmoji},
			{Emoji: denyEmoji},
		},
	}

	if _, err := c.Client.InsertOne(context.Background(), request); err != nil {
		return nil, fmt.Errorf("error creating nickname request: %v", err)
	}
	return request, nil
}

func (c *MongoClient) RemoveNicknameRequest(messageId string) (int, error) {
	if messageId == "" {
		return 0, fmt.Errorf("messageId cannot be empty")
	}

	res, err := c.Client.DeleteOne(context.Background(), bson.M{"user_message_id": messageId})
	if err != nil {
		return 0, fmt.Errorf("error deleting nickname request with messageId: %s, %v", messageId, err)
	}

	return int(res.DeletedCount), nil
}

func (c *MongoClient) SetupNicknameRequestCollector(s *discordgo.Session, message *discordgo.Message, nickname string) error {
    fmt.Printf("Setting up reaction collector for nickname request: %s on message ID: %s\n", nickname, message.ID)

    // Handler function for reactions
    handler := func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
        // Only process reactions for the target message
        if r.MessageID != message.ID {
            return
        }

        // Get guild and member
        guildID := r.GuildID
        if guildID == "" {
            return
        }
        member, err := s.GuildMember(guildID, r.UserID)
        if err != nil || member == nil || r.UserID == s.State.User.ID {
            return
        }

        // Check staff role
        staffRoleID := os.Getenv("STAFF_ROLE_ID") // Or use your STAFF_ROLE_IDS[0]
        hasStaffRole := false
        for _, roleID := range member.Roles {
            if roleID == staffRoleID {
                hasStaffRole = true
                break
            }
        }
        if !hasStaffRole {
            return
        }

        // Check emoji
        if r.Emoji.ID != approvedEmojiID && r.Emoji.ID != deniedEmojiID {
            return
        }

        // Fetch nickname request from DB
        request, err := repository.NicknameRequestService.DBClient.GetNicknameRequestByStaffMessageID(message.ID)
        if err != nil || request == nil {
            return
        }

        // Get the user to change
        userToChange, err := s.GuildMember(guildID, request.UserID)
        if err != nil || userToChange == nil {
            return
        }

        // Approve or deny
        switch r.Emoji.ID {
					case approvedEmojiID:
            // Change nickname
            err := s.GuildMemberNickname(guildID, request.UserID, nickname)
            if err == nil {
                fmt.Printf("Changed nickname for %s to %s\n", userToChange.User.Username, nickname)
            }
            // React to user message
            s.MessageReactionAdd(request.UserChannelID, request.UserMessageID, approveEmoji)
            repository.NicknameRequestService.DBClient.RemoveNicknameRequest(message.ID)
        	case deniedEmojiID:
            repository.NicknameRequestService.DBClient.RemoveNicknameRequest(message.ID)
            s.MessageReactionAdd(request.UserChannelID, request.UserMessageID, denyEmoji)
            fmt.Printf("Nickname request for %s to %s is denied\n", userToChange.User.Username, nickname)
					case "":
						// Handle case where emoji is not recognized
						fmt.Printf("Unrecognized emoji reaction: %s\n", r.Emoji.Name)
        }
    }

    // Add the handler
    s.AddHandler(handler)
    return nil
}

func (c *MongoClient) GetNicknameRequestByStaffMessageID(messageID string) (*model.NicknameRequest, error) {
	if messageID == "" {
		return nil, fmt.Errorf("messageID cannot be empty")
	}

	var request model.NicknameRequest
	err := c.Client.FindOne(context.Background(), bson.M{"staff_message_id": messageID}).Decode(&request)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No request found
		}
		return nil, fmt.Errorf("error finding nickname request by staff message ID: %v", err)
	}
	return &request, nil
}

func (c *MongoClient) InitializeNicknameRequests(s *discordgo.Session) error {
	guildID := os.Getenv("GUILD_ID")
	guild, err := s.Guild(guildID)
	if err != nil || guild == nil {
		return fmt.Errorf("guild not found: %v", err)
	}

	nicknameRequests, err := repository.NicknameRequestService.DBClient.GetAllNicknameRequests()
	if err != nil {
		return fmt.Errorf("error fetching nickname requests: %v", err)
	}

	if len(nicknameRequests) == 0 {
		fmt.Println("No nickname requests found.")
		return nil
	}

	channels := guild.Channels
	if len(channels) == 0 {
		return fmt.Errorf("no channels found in guild: %s", guildID)
	}

	for _, channel := range channels {

		filteredNicknameRequests := make([]*model.NicknameRequest, 0)
    for _, nicknameRequest := range nicknameRequests {
			if nicknameRequest.StaffChannelID == channel.ID {
				filteredNicknameRequests = append(filteredNicknameRequests, nicknameRequest)
			}
    }

		for _, filteredRequest := range filteredNicknameRequests {
			message, err := s.ChannelMessage(channel.ID, filteredRequest.StaffMessageID)
			if err != nil || message == nil {
				fmt.Printf("Error fetching message with ID %s from channel %s: %v\n", filteredRequest.StaffMessageID, channel.ID, err)
				continue
			}

			err = repository.NicknameRequestService.DBClient.SetupNicknameRequestCollector(s, message, filteredRequest.Nickname)
			if err != nil {
				fmt.Printf("Error setting up nickname request collector for message ID %s: %v\n", filteredRequest.StaffMessageID, err)
				continue
			}
			fmt.Printf("Initialized react %s for message %s in channel %s\n", filteredRequest.Nickname, filteredRequest.StaffMessageID, channel.ID)
		}
	}
	fmt.Println("Nickname requests initialized successfully.")
	return nil
}