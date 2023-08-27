package entities

type Message struct {
	Id   int    `db:"id"`
	Text string `db:"msg"`
}
