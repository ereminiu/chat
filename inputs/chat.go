package inputs

type Chat struct {
	Id    int      `json:"id"`
	Name  string   `json:"name"`
	Users []string `json:"users"`
}
