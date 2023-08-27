package entities

type Chat struct {
	Id   int    `db:"id"`
	Name string `json:"name"`
}
