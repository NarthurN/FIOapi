package user

type User struct {
	ID          int
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	Sex         string
	Nationality string
}

type UserJSON struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}
