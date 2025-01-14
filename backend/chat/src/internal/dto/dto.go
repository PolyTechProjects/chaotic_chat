package dto

import (
	"github.com/PolyTechProjects/chaotic_chat/chat/src/internal/models"
	"github.com/google/uuid"
)

type CreateChatRequest struct {
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	CreatorId       uuid.UUID   `json:"creator_id"`
	ParticipantsIds []uuid.UUID `json:"participants_ids"`
	IsChannel       bool        `json:"is_channel"`
	IsClosed        bool        `json:"is_closed"`
}

type UpdateChatRequest struct {
	ChatId      uuid.UUID `json:"chat_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ProfilePic  string    `json:"profile_pic"`
}

type UpdateUsersRequest struct {
	ChatId  uuid.UUID   `json:"chat_id"`
	UserIds []uuid.UUID `json:"users_ids"`
}

type UpdateAdminsRequest struct {
	ChatId    uuid.UUID   `json:"chat_id"`
	AdminsIds []uuid.UUID `json:"admins_ids"`
}

type UpdateReadersRequest struct {
	ChatId     uuid.UUID   `json:"chat_id"`
	ReadersIds []uuid.UUID `json:"readers_ids"`
}

type GetChatResponse struct {
	Chat    *models.Chat
	Users   []string
	Admins  []string
	Readers []string
}
