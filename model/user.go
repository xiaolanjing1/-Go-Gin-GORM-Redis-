package model

import "time"

type User struct {
	Id        int
	Name      string `json:"name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Avatar    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
