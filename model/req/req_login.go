package req
type ReqLogIn struct{
	//FullName string `json:"fullName,omitempty" validate:"required"`
	Email string	`json:"email,omitempty" validate:"required,email"`
	Password string	`json:"password,omitempty" validate:"required"`
}