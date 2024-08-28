package model

import "time"

type GithubRepo struct {
	Name        string  `json:"name,omitempty" db:"name, omitempty"`
	Description string  `json:"description,omitempty" db:"description, omitempty"`
	Url         string	`json:"url,omitempty" db:"url, omitempty"`
	Color       string	`json:"color,omitempty" db:"color, omitempty"`
	Lang        string	`json:"lang,omitempty" db:"lang, omitempty"`
	Fork        string	`json:"fork,omitempty" db:"fork, omitempty"`
	Stars       string	`json:"stars,omitempty" db:"stars, omitempty"`
	Stars_today string	`json:"starsToday,omitempty" db:"stars_today, omitempty"`
	Build_by    string	`json:"buildBy,omitempty" db:"build_by, omitempty"`
	Bookmarked string 	`json:"bookmarked,omitempty" db:"bookmarked, omitempty"`
	Created_at  time.Time	`json:"-" db:"created_at, omitempty"`
	Updated_at time.Time	`json:"-" db:"updated_at, omitempty"`
}