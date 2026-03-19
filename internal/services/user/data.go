package userService

import "pantho/golang/internal/models"

type ResponseCtx struct {
	Name string        `json:"name"`
	Data []models.User `json:"data"`
}