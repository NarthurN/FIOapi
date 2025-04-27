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

type UsersResponse struct {
	Users      []User `json:"users"`
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
}

type UserFilter struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	AgeFrom     int    `json:"age_from"`
	AgeTo       int    `json:"age_to"`
	Sex         string `json:"sex"`
	Nationality string `json:"nationality"`
}

type Pagination struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}
