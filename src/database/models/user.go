package models

type User struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	Name string `json:"name"`
}
