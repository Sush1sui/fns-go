package repository

import "github.com/Sush1sui/fns-go/internal/model"

var StickyChannels = map[string]struct{}{}

type StickyInterface interface {
	InitializeStickyChannels() error
	GetAllStickyChannels() ([]*model.StickyChannel, error)
	GetStickyChannel(string) *model.StickyChannel
	CreateStickyChannel(string, string) (string, error)
	GetRecentPostMessageId(string) (string, error)
	UpdateStickyMessageId(string, string, string) (*model.StickyChannel, error)
	SetStickyMessageId(string, string) (*model.StickyChannel, error)
	DeleteStickyChannel(string) (int, error)
	DeleteAllStickyChannels() (int, error)
}

type StickyServiceType struct {
	DBClient StickyInterface
}

var StickyService StickyServiceType