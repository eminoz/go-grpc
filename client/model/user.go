package model

type User struct {
	Name     string `json:"name"`
	Email    string `validate:"required,email,omitempty"`
	Password string `validate:"required,gte=7,lte=130,omitempty"`
	Role     string `json:"role"`
}
