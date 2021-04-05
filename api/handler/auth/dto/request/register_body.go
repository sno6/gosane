package request

type RegisterBody struct {
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}
