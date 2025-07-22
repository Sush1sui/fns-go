package repository

import (
	"github.com/Sush1sui/fns-go/internal/model"
	"github.com/bwmarrin/discordgo"
)

type NicknameRequestInterface interface {
	GetAllNicknameRequests() ([]*model.NicknameRequest, error)
	CreateNicknameRequest(val, userId, userMessageId, userChannelId, staffChannelId, staffMessageId string) (*model.NicknameRequest, error)
	RemoveNicknameRequest(messageId string) (int, error)
	SetupNicknameRequestCollector(session *discordgo.Session, message *discordgo.Message, nickname string) error
	GetNicknameRequestByStaffMessageID(messageID string) (*model.NicknameRequest, error)
	InitializeNicknameRequests(s *discordgo.Session) error
}

type NicknameRequestServiceType struct {
	DBClient NicknameRequestInterface
}

var NicknameRequestService NicknameRequestServiceType