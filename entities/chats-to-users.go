package entities

type ChatsToUsers struct {
	ChatId   int    `db:"chat_id"`
	ChatName string `db:"chat_name"`
	UserName string `db:"user_name"`
	UserId   int    `db:"user_id"`
}
