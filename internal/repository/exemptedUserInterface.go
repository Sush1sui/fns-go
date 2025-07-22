package repository

import "github.com/Sush1sui/fns-go/internal/model"

type ExemptedUserInterface interface {
	ExemptUserVanity(string) (bool, error)
	RemoveExemptedUser(string) (int, error)
	GetAllExemptedUsers() ([]*model.ExemptedUser, error)
}

type ExemptedServiceType struct {
	DBClient ExemptedUserInterface
}

var ExemptedService ExemptedServiceType