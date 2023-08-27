package inputs

type User struct {
	Id   int    `db:"id" json:"id"`
	Name string `json:"user_name"`
}
