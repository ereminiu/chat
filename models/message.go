package models

type Message struct {
	Id     int    `db:"id" json:"id"`
	Chat   int    `db:"chat" json:"chat"`
	Author int    `db:"author" json:"author"`
	Text   string `db:"text" json:"text"`
}

type DebugMessage struct {
	Id     int    `db:"id" json:"id"`
	Author string `db:"author" json:"author"`
	Text   string `db:"text" json:"text"`
}
