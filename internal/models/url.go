package models

type Url struct {
	Id    string `json:"id" db:"id"`
	Url   string `json:"url" db:"url"`
	Alias string `json:"alias" db:"alias"`
}
