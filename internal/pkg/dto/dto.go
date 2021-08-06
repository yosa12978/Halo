package dto

type UserInfo struct {
	User_id string
	Roles   []string
}

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

type CreateUser struct {
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	BirthDate int64    `json:"birthDate"`
	City      string   `json:"city"`
	Country   string   `json:"country"`
	Skills    []string `json:"skills"`
}

type UpdateUser struct {
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	BirthDate int64    `json:"birthDate"`
	City      string   `json:"city"`
	Country   string   `json:"country"`
	Skills    []string `json:"skills"`
	Status    int8     `json:"status"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateOffer struct {
	Title        string   `json:"title"`
	Body         string   `json:"body"`
	Tags         []string `json:"tags"`
	Company      string   `json:"company"`
	City         string   `json:"city"`
	ContactEmail string   `json:"contactEmail"`
	Website      string   `json:"website"`
}

type UpdateOffer struct {
	Title        string   `json:"title"`
	Body         string   `json:"body"`
	Tags         []string `json:"tags"`
	City         string   `json:"city"`
	ContactEmail string   `json:"contactEmail"`
}
