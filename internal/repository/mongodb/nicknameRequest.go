package mongodb

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/Sush1sui/fns-go/internal/common"
	"github.com/Sush1sui/fns-go/internal/model"
	"github.com/Sush1sui/fns-go/internal/repository"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var approveEmoji string
var denyEmoji string
var approvedEmojiID string
var deniedEmojiID string

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

	res := c.Client.FindOneAndDelete(context.Background(), bson.M{"userMessageId": messageId})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return 0, nil // No document found to delete
		}
		return 0, fmt.Errorf("error deleting nickname request with messageId: %s, %v", messageId, res.Err())
	}
	// Since FindOneAndDelete does not return DeletedCount, assume 1 if no error
	return 1, nil
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
        hasStaffRole := slices.Contains(member.Roles, common.StaffRoleIDs[0])
        if !hasStaffRole {
            return
        }

        // Check emoji
        if r.Emoji.ID != approvedEmojiID && r.Emoji.ID != deniedEmojiID {
            fmt.Println("Unrecognized emoji reaction:", r.Emoji.ID)
            return
        }

        // Fetch nickname request from DB
        request, err := repository.NicknameRequestService.DBClient.GetNicknameRequestByStaffMessageID(message.ID)
        if err != nil || request == nil {
            fmt.Println("Error fetching nickname request:", err)
            return
        }

        // Get the user to change
        userToChange, err := s.GuildMember(guildID, request.UserID)
        if err != nil || userToChange == nil {
						fmt.Println("Error fetching user to change nickname:", err)
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
            err = s.MessageReactionAdd(request.UserChannelID, request.UserMessageID, "Check_White_FNS:1310274014102687854")
            if err != nil {
                fmt.Println("Error reacting to user message:", err)
            }
            count, err := repository.NicknameRequestService.DBClient.RemoveNicknameRequest(request.UserMessageID)
						if err != nil {
								fmt.Println("Error removing nickname request:", err)
						} else if count == 0 {
								fmt.Println("No nickname request was removed (count == 0)")
						}
        	case deniedEmojiID:
            count, err := repository.NicknameRequestService.DBClient.RemoveNicknameRequest(request.UserMessageID)
            if err != nil {
								fmt.Println("Error removing nickname request:", err)
						} else if count == 0 {
								fmt.Println("No nickname request was removed (count == 0)")
						}
            err = s.MessageReactionAdd(request.UserChannelID, request.UserMessageID, "No:1310633209519669290")
            if err != nil {
                fmt.Println("Error reacting to user message:", err)
            }
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
	err := c.Client.FindOne(context.Background(), bson.M{"staffMessageId": messageID}).Decode(&request)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No request found
		}
		return nil, fmt.Errorf("error finding nickname request by staff message ID: %v", err)
	}
	return &request, nil
}

func (c *MongoClient) InitializeNicknameRequests(s *discordgo.Session) error {
	approveEmoji = os.Getenv("APPROVE_EMOJI")
	denyEmoji = os.Getenv("DENY_EMOJI")
	approvedEmojiID = os.Getenv("APPROVED_EMOJI_ID")
	deniedEmojiID = os.Getenv("DENIED_EMOJI_ID")
	
	guildID := os.Getenv("GUILD_ID")
	guild, err := s.State.Guild(guildID)
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