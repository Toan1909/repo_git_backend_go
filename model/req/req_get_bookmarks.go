package req

type ReqBookmark struct {
	NameRepo string `json:"repoName,omitempty" validate:"required"`
}
