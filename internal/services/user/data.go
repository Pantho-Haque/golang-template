package userService

import "pantho/golang/internal/models"

type ResponseCtx struct {
	Time string        `json:"time"`
	Data []models.User `json:"data"`
}