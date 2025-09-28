package model

import (
	storage "github.com/LTSlw/QiniuIIPorject/backend/pkg/storge"
)

type ChatResponse struct {
	SessionId int             `json:"session_id"`
	Message   storage.Message `json:"message"`
}
