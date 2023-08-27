package repository

import (
	"fmt"
	"log"
	"strings"

	"github.com/ereminiu/chat/entities"
	"github.com/ereminiu/chat/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository() (*Repository, error) {
	dataSource := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		"ys-user", "qwerty", "localhost", 5432, "ys-db",
	)
	conn, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		return nil, errors.Wrap(err, "sqlx connect")
	}
	if err := conn.Ping(); err != nil {
		return nil, errors.Wrap(err, "ping failed")
	}
	return &Repository{db: conn}, nil
}

func (r *Repository) AddUser(user *entities.User) error {
	var userId int
	err := r.db.Get(
		&userId,
		`
			INSERT INTO users (name) 
			VALUES ($1)
			RETURNING id
		`,
		user.Name,
	)
	user.Id = userId
	if err != nil {
		return errors.Wrap(err, "insert user")
	}
	return nil
}

func (r *Repository) AddChat(chat *entities.Chat) error {
	var chatId int
	err := r.db.Get(
		&chatId,
		`
			INSERT INTO chats (name)
			VALUES ($1)
			RETURNING id
		`,
		chat.Name,
	)
	chat.Id = chatId
	if err != nil {
		return errors.Wrap(err, "insert chat")
	}
	return nil
}

func (r *Repository) AddUserToChat(chat *entities.Chat, userId int) error {
	tx := r.db.MustBegin()
	tx.MustExec(
		`
			INSERT INTO users_to_chats (user_id, chat_id)
			VALUES ($1, $2)
		`,
		userId, chat.Id,
	)
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "insert user to chat")
	}
	return nil
}

func (r *Repository) GetUsersChats(userId int) ([]entities.ChatsToUsers, error) {
	UsersChats := make([]string, 0)
	err := r.db.Select(&UsersChats,
		`
			SELECT c.name 
			FROM users u
			LEFT JOIN users_to_chats uc
			ON uc.user_id = u.id 
			JOIN chats c
			ON uc.chat_id = c.id
			WHERE u.id = ($1)
		`,
		userId,
	)
	if len(UsersChats) == 0 {
		return make([]entities.ChatsToUsers, 0), nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "select users chats")
	}
	inParams := make([]string, 0, len(UsersChats))
	args := make([]interface{}, 0, len(UsersChats))
	for i, chat := range UsersChats {
		inParams = append(inParams, fmt.Sprintf("$%d", i+1))
		args = append(args, interface{}(chat))
	}
	query := fmt.Sprintf(
		`
			SELECT c.id as "chat_id", c.name as "chat_name", u.name as "user_name", u.id as "user_id" 
			FROM chats c 
			LEFT JOIN users_to_chats uc 
			ON uc.chat_id = c.id
			LEFT JOIN users u
			ON uc.user_id = u.id
			WHERE c.name in (%s)
			ORDER BY c.name
		`,
		strings.Join(inParams, ","),
	)
	rawData := make([]entities.ChatsToUsers, 0)
	if err := r.db.Select(&rawData, query, args...); err != nil {
		return nil, errors.Wrap(err, "select data of the chats")
	}
	return rawData, nil
}

func (r *Repository) AddMessage(message *entities.Message) error {
	var messageId int
	err := r.db.Get(
		&messageId,
		`
			INSERT INTO messages (msg)
			VALUES ($1)
			RETURNING id
		`,
		message.Text,
	)
	if err != nil {
		return errors.Wrap(err, "insert message")
	}
	message.Id = messageId
	return nil
}

func (r *Repository) AddMessageToChat(message *entities.Message, chatId int) error {
	tx := r.db.MustBegin()
	tx.MustExec(
		`
			INSERT INTO messages_to_chats (message_id, chat_id)
			VALUES ($1, $2)
		`,
		message.Id, chatId,
	)
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "insert message to chat")
	}
	return nil
}

func (r *Repository) AddMessageToUser(message *entities.Message, userId int) error {
	tx := r.db.MustBegin()
	tx.MustExec(
		`
			INSERT INTO messages_to_users (user_id, message_id)
			VALUES ($1, $2)
		`,
		userId, message.Id,
	)
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "insert message to user")
	}
	return nil
}

func (r *Repository) GetMessagesByChat(chatId int) ([]models.Message, error) {
	rawData := make([]models.Message, 0)
	err := r.db.Select(
		&rawData,
		`
			SELECT m.id as "id", c.id as "chat", u.id as "author", m.msg as "text"
			FROM users u
			JOIN messages_to_users mu
			ON mu.user_id = u.id
			JOIN messages m
			ON m.id = mu.message_id
			JOIN messages_to_chats mc
			ON mc.message_id = m.id
			JOIN chats c
			ON mc.chat_id = c.id
			WHERE c.id = $1
		`,
		chatId,
	)
	if err != nil {
		return nil, errors.Wrap(err, "select messages from chat")
	}
	return rawData, nil
}

func (r *Repository) DebugGetMessagesByChat(chatId int) ([]models.DebugMessage, error) {
	rawData := make([]models.DebugMessage, 0)
	err := r.db.Select(
		&rawData,
		`
			SELECT m.id as "id", u.name as "author", m.msg as "text"
			FROM users u
			JOIN messages_to_users mu
			ON mu.user_id = u.id
			JOIN messages m
			ON m.id = mu.message_id
			JOIN messages_to_chats mc
			ON mc.message_id = m.id
			JOIN chats c
			ON mc.chat_id = c.id
			WHERE c.id = $1
		`,
		chatId,
	)
	if err != nil {
		return nil, errors.Wrap(err, "select messages from chat")
	}
	return rawData, nil
}

func (r *Repository) UserIdExists(userId int) bool {
	var exists bool
	err := r.db.Get(
		&exists,
		`
			select exists(select 1 from users where id = $1 limit 1)
		`,
		userId,
	)
	if err != nil {
		log.Fatal(err)
	}
	return exists
}

func (r *Repository) UserNameExists(userName string) bool {
	var exists bool
	err := r.db.Get(
		&exists,
		`
			select exists(select 1 from users where name = $1 limit 1)
		`,
		userName,
	)
	if err != nil {
		log.Fatal(err)
	}
	return exists
}

func (r *Repository) ChatIdExists(chatId int) bool {
	var exists bool
	err := r.db.Get(
		&exists,
		`
			select exists(select 1 from chats where id = $1 limit 1)
		`,
		chatId,
	)
	if err != nil {
		log.Fatal(err)
	}
	return exists
}

func (r *Repository) ChatNameExists(ChatName string) bool {
	var exists bool
	err := r.db.Get(
		&exists,
		`
			select exists(select 1 from chats where name = $1 limit 1)
		`,
		ChatName,
	)
	if err != nil {
		log.Fatal(err)
	}
	return exists
}
