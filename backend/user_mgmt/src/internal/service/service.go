package service

import (
	"fmt"
	"log/slog"

	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/internal/repository"
	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/models"
	"github.com/google/uuid"
)

type UserMgmtService struct {
	Repository *repository.UserMgmtRepository
}

func New(repository *repository.UserMgmtRepository) *UserMgmtService {
	return &UserMgmtService{Repository: repository}
}

func (s *UserMgmtService) CreateUser(userId uuid.UUID, name string) (*models.User, error) {
	user := models.New(userId, name)
	err := s.Repository.InsertUser(user)
	slog.Info(fmt.Sprintf("User %v created", user.Id))
	return user, err
}

func (s *UserMgmtService) UpdateUser(userId uuid.UUID, name string, urlTag string, description string) (*models.User, error) {
	user, err := s.GetUser(userId)
	if err != nil {
		return nil, err
	}
	user.Name = name
	user.UrlTag = urlTag
	user.Description = description
	err = s.Repository.UpdateUser(user)
	return user, err
}

func (s *UserMgmtService) DeleteUser(userId uuid.UUID) error {
	return s.Repository.DeleteUser(userId)
}

func (s *UserMgmtService) GetUser(userId uuid.UUID) (*models.User, error) {
	user, err := s.Repository.GetUser(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserMgmtService) UpdateAvatar(userId uuid.UUID, newFileId string) (*models.User, error) {
	user, err := s.GetUser(userId)
	if err != nil {
		return nil, err
	}
	user.ProfilePic = newFileId
	err = s.Repository.UpdateUser(user)
	return user, err
}
