package models

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Role     string    `json:"role"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type CreateUser struct {
	FullName string `json:"full_name" validate:"required"`
	Role     string `json:"role" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}

type UserFilter struct {
	Search string
	Limit  int
	Offset int
}

type UpdateUser struct {
	Id       uuid.UUID `json:"-"`
	FullName string    `json:"full_name" validate:"required"`
	Role     string    `json:"role" validate:"required"`
	Username string    `json:"username" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required"`
}
