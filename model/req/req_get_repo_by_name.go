package req
type ReqGetRepoByName struct{
	Name string	`json:"name,omitempty" validate:"required, name"`
}