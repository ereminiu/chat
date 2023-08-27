package models

import "github.com/ereminiu/chat/entities"

type Chat struct {
	Id    int             `json:"id"`
	Name  string          `json:"name"`
	Users []entities.User `json:"users"`
}
