package handler

type FullName struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

type UrlParams struct {
	Gender      string
	Nationality string
	Page        int
	Size        int
}
