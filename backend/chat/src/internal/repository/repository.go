package repository

import (
	"github.com/PolyTechProjects/chaotic_chat/chat/src/internal/models"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) FindById(chatId uuid.UUID) (*models.Chat, error) {
	var chat models.Chat
	err := r.db.Where("id = ?", chatId).First(&chat).Error
	return &chat, err
}

func (r *ChatRepository) FindByJoinLink(joinLink string) (*models.Chat, error) {
	var chat models.Chat
	err := r.db.Where("join_link = ?", joinLink).First(&chat).Error
	return &chat, err
}

func (r *ChatRepository) SaveChat(chat *models.Chat) error {
	return r.db.Create(chat).Error
}

func (r *ChatRepository) DeleteChat(chatId uuid.UUID) error {
	return r.db.Where("id = ?", chatId).Delete(&models.Chat{}).Error
}

func (r *ChatRepository) UpdateChat(chat *models.Chat) error {
	return r.db.Save(chat).Error
}

func (r *ChatRepository) AddUserToChat(user *models.UserChat) error {
	return r.db.Create(&user).Error
}

func (r *ChatRepository) RemoveUserFromChat(user *models.UserChat) error {
	return r.db.Where("chat_id = ? AND user_id = ?", user.ChatId, user.UserId).Delete(&models.UserChat{}).Error
}

func (r *ChatRepository) UpdateUserChat(user *models.UserChat) error {
	return r.db.Save(user).Error
}

func (r *ChatRepository) GetChatUsers(chatID uuid.UUID) ([]models.UserChat, error) {
	var userChats []models.UserChat
	err := r.db.Where("chat_id = ?", chatID).Find(&userChats).Error
	if err != nil {
		return nil, err
	}

	if len(userChats) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return userChats, nil
}
