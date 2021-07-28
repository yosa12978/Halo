package dto

type CreateCompany struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Website  string `json:"website"`
	Contact  string `json:"contact"`
}

type UpdateCompany struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Website string `json:"website"`
	Contact string `json:"contact"`
}
