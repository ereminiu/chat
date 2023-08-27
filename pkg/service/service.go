package service

import (
	"log"

	"github.com/pkg/errors"

	"github.com/ereminiu/chat/entities"
	"github.com/ereminiu/chat/models"
	"github.com/ereminiu/chat/pkg/repository"
)

type Serivce struct {
	repos *repository.Repository
}

func NewService(r *repository.Repository) (*Serivce, error) {
	return &Serivce{repos: r}, nil
}

// Add user and set ID
func (s *Serivce) AddUser(user *entities.User) error {
	return s.repos.AddUser(user)
}

func (s *Serivce) CreateChat(chat *entities.Chat, userIds []int) error {
	if err := s.AddChat(chat); err != nil {
		return errors.Wrap(err, "chat")
	}
	for _, userId := range userIds {
		if err := s.AddUserToChat(chat, userId); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

// Add chat and set ID
func (s *Serivce) AddChat(chat *entities.Chat) error {
	return s.repos.AddChat(chat)
}

func (s *Serivce) AddUserToChat(chat *entities.Chat, userId int) error {
	return s.repos.AddUserToChat(chat, userId)
}

func (s *Serivce) GetUsersChats(userId int) ([]models.Chat, error) {
	rawData, err := s.repos.GetUsersChats(userId)
	if err != nil {
		return nil, err
	}
	prevChatName := "@"
	chats := make([]models.Chat, 0)
	for _, cur := range rawData {
		if cur.ChatName != prevChatName {
			chats = append(chats, models.Chat{Id: cur.ChatId, Name: cur.ChatName, Users: make([]entities.User, 0)})
		}
		last := len(chats) - 1
		chats[last].Users = append(chats[last].Users, entities.User{Id: cur.UserId, Name: cur.UserName})
		prevChatName = cur.ChatName
	}
	return chats, nil
}

func (s *Serivce) AddMessage(message *entities.Message) error {
	return s.repos.AddMessage(message)
}

func (s *Serivce) AddMessageToChat(message *entities.Message, chatId int) error {
	return s.repos.AddMessageToChat(message, chatId)
}

func (s *Serivce) AddMessageToUser(message *entities.Message, userId int) error {
	return s.repos.AddMessageToUser(message, userId)
}

func (s *Serivce) GetMessagesByChat(chatId int) ([]models.Message, error) {
	return s.repos.GetMessagesByChat(chatId)
}

func (s *Serivce) DebugGetMessagesByChat(chatId int) ([]models.DebugMessage, error) {
	return s.repos.DebugGetMessagesByChat(chatId)
}

func (s *Serivce) UserIdExists(userId int) bool {
	return s.repos.UserIdExists(userId)
}

func (s *Serivce) UserNameExists(userName string) bool {
	return s.repos.UserNameExists(userName)
}

func (s *Serivce) ChatIdExists(chatId int) bool {
	return s.repos.ChatIdExists(chatId)
}

func (s *Serivce) ChatNameExists(chatName string) bool {
	return s.repos.ChatNameExists(chatName)
}
