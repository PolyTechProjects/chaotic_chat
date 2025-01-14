package service

import (
	"errors"
	"log/slog"

	"github.com/PolyTechProjects/chaotic_chat/chat/src/internal/dto"
	"github.com/PolyTechProjects/chaotic_chat/chat/src/internal/models"
	"github.com/PolyTechProjects/chaotic_chat/chat/src/internal/repository"
	"github.com/google/uuid"
)

type ChatManagementService struct {
	repo repository.ChatRepository
}

func New(repo repository.ChatRepository) *ChatManagementService {
	return &ChatManagementService{
		repo: repo,
	}
}

func (s *ChatManagementService) GetChat(chatId uuid.UUID) (*dto.GetChatResponse, error) {
	chat, err := s.repo.FindById(chatId)
	if err != nil {
		slog.Error("Failed to get chat", "error", err.Error())
		return nil, err
	}
	userChats, err := s.repo.GetChatUsers(chatId)
	if err != nil {
		slog.Error("Failed to get chat users", "error", err.Error())
		return nil, err
	}
	var users []string
	var admins []string
	var readers []string
	for _, userChat := range userChats {
		users = append(users, userChat.UserId.String())
		if userChat.IsAdmin {
			admins = append(admins, userChat.UserId.String())
		} else if userChat.ReadOnly {
			readers = append(readers, userChat.UserId.String())
		}
	}
	getResp := &dto.GetChatResponse{
		Chat:    chat,
		Users:   users,
		Admins:  admins,
		Readers: readers,
	}
	return getResp, nil
}

func (s *ChatManagementService) CreateChat(req *dto.CreateChatRequest) (*dto.GetChatResponse, error) {
	slog.Info("CreateChat called", "name", req.Name, "description", req.Description, "creatorID", req.CreatorId)
	chat := &models.Chat{
		Id:          uuid.New(),
		Name:        req.Name,
		CreatorId:   req.CreatorId,
		IsChannel:   req.IsChannel,
		Description: req.Description,
		ProfilePic:  "",
	}
	err := s.repo.SaveChat(chat)
	if err != nil {
		slog.Error("Failed to create chat", "error", err.Error())
		return nil, err
	}
	creatorChat := models.UserChat{
		ChatId:   chat.Id,
		UserId:   req.CreatorId,
		IsAdmin:  true,
		ReadOnly: false,
	}
	err = s.repo.AddUserToChat(&creatorChat)
	if err != nil {
		slog.Error("Failed to add user to chat", "error", err.Error())
		return nil, err
	}
	readers := []string{}
	users := []string{}
	admins := []string{req.CreatorId.String()}
	for _, participantId := range req.ParticipantsIds {
		userChat := models.UserChat{
			ChatId:   chat.Id,
			UserId:   participantId,
			IsAdmin:  false,
			ReadOnly: req.IsChannel,
		}
		err = s.repo.AddUserToChat(&userChat)
		if err != nil {
			slog.Error("Failed to add user to chat", "error", err.Error())
			return nil, err
		}
		readers = append(readers, participantId.String())
	}
	return &dto.GetChatResponse{Chat: chat, Admins: admins, Readers: readers, Users: users}, nil
}

func (s *ChatManagementService) DeleteChat(chatId uuid.UUID) error {
	err := s.repo.DeleteChat(chatId)
	if err != nil {
		slog.Error("Failed to delete chat", "error", err.Error())
		return err
	}
	return nil
}

func (s *ChatManagementService) UpdateChat(req *dto.UpdateChatRequest) error {
	chat, err := s.repo.FindById(req.ChatId)
	if err != nil {
		slog.Error("Failed to find chat", "error", err.Error())
		return err
	}
	chat.Name = req.Name
	chat.Description = req.Description
	chat.ProfilePic = req.ProfilePic
	err = s.repo.UpdateChat(chat)
	if err != nil {
		slog.Error("Failed to update chat", "error", err.Error())
		return err
	}
	return nil
}

func (s *ChatManagementService) JoinChat(joinLink string, userId uuid.UUID) error {
	chat, err := s.repo.FindByJoinLink(joinLink)
	if err != nil {
		slog.Error("Failed to find chat by join link", "error", err.Error())
		return err
	}
	userChats, err := s.repo.GetChatUsers(chat.Id)
	if err != nil {
		slog.Error("Failed to get chat users", "error", err.Error())
		return err
	}
	if len(userChats) > 19 && !chat.IsChannel {
		return errors.New("chat is full")
	}
	userChat := &models.UserChat{
		ChatId:   chat.Id,
		UserId:   userId,
		ReadOnly: chat.IsChannel,
		IsAdmin:  false,
	}
	err = s.repo.AddUserToChat(userChat)
	if err != nil {
		slog.Error("Failed to add user to chat", "error", err.Error())
		return err
	}
	return nil
}

func (s *ChatManagementService) DeleteUsers(chatId uuid.UUID, userIds []uuid.UUID) error {
	chat, err := s.repo.FindById(chatId)
	if err != nil {
		slog.Error("Failed to find chat", "error", err.Error())
		return err
	}
	for _, userId := range userIds {
		if chat.CreatorId == userId {
			slog.Error("Cannot delete creator")
			return errors.New("cannot delete creator")
		}
		userChats, err := s.repo.GetChatUsers(chatId)
		if err != nil {
			slog.Error("Failed to get chat users", "error", err.Error())
			return err
		}
		if len(userChats)-len(userIds) < 2 && !chat.IsChannel {
			return s.DeleteChat(chatId)
		}
		for _, userChat := range userChats {
			if userChat.UserId == userId {
				err = s.repo.RemoveUserFromChat(&userChat)
				if err != nil {
					slog.Error("Failed to remove user from chat", "error", err.Error())
					return err
				}
			}
		}
	}
	return nil
}

func (s *ChatManagementService) AddUsers(chatId uuid.UUID, userIds []uuid.UUID) error {
	chat, err := s.repo.FindById(chatId)
	if err != nil {
		slog.Error("Failed to find chat", "error", err.Error())
		return err
	}
	userChats, err := s.repo.GetChatUsers(chat.Id)
	if err != nil {
		slog.Error("Failed to get chat users", "error", err.Error())
		return err
	}
	if len(userChats) > 19 && !chat.IsChannel {
		return errors.New("chat is full")
	}
	for _, userId := range userIds {
		userChat := models.UserChat{
			ChatId:   chatId,
			UserId:   userId,
			ReadOnly: chat.IsChannel,
			IsAdmin:  false,
		}
		err := s.repo.AddUserToChat(&userChat)
		if err != nil {
			slog.Error("Failed to add user to chat", "error", err.Error())
			return err
		}
	}
	return nil
}

func (s *ChatManagementService) MakeReadersUsers(chatId uuid.UUID, readersIds []uuid.UUID) error {
	userChats, err := s.repo.GetChatUsers(chatId)
	if err != nil {
		slog.Error("Failed to get chat users", "error", err.Error())
		return err
	}
	for _, readerId := range readersIds {
		for _, userChat := range userChats {
			if userChat.UserId == readerId {
				userChat.ReadOnly = false
				err = s.repo.UpdateUserChat(&userChat)
				if err != nil {
					slog.Error("Failed to update user chat", "error", err.Error())
					return err
				}
			}
		}
	}
	return nil
}

func (s *ChatManagementService) MakeUsersReaders(chatId uuid.UUID, usersIds []uuid.UUID) error {
	userChats, err := s.repo.GetChatUsers(chatId)
	if err != nil {
		slog.Error("Failed to get chat users", "error", err.Error())
		return err
	}
	for _, userId := range usersIds {
		for _, userChat := range userChats {
			if userChat.UserId == userId {
				userChat.ReadOnly = true
				userChat.IsAdmin = false
				err = s.repo.UpdateUserChat(&userChat)
				if err != nil {
					slog.Error("Failed to update user chat", "error", err.Error())
					return err
				}
			}
		}
	}
	return nil
}

func (s *ChatManagementService) AddAdmin(chatId uuid.UUID, adminsIds []uuid.UUID) error {
	userChats, err := s.repo.GetChatUsers(chatId)
	if err != nil {
		slog.Error("Failed to get chat users", "error", err.Error())
		return err
	}
	for _, adminId := range adminsIds {
		for _, userChat := range userChats {
			if userChat.UserId == adminId {
				userChat.IsAdmin = true
				err = s.repo.UpdateUserChat(&userChat)
				if err != nil {
					slog.Error("Failed to update user chat", "error", err.Error())
					return err
				}
			}
		}
	}
	return nil
}

func (s *ChatManagementService) DeleteAdmin(chatId uuid.UUID, adminsIds []uuid.UUID) error {
	userChats, err := s.repo.GetChatUsers(chatId)
	if err != nil {
		slog.Error("Failed to get chat users", "error", err.Error())
		return err
	}
	for _, adminId := range adminsIds {
		for _, userChat := range userChats {
			if userChat.UserId == adminId {
				userChat.IsAdmin = false
				err = s.repo.UpdateUserChat(&userChat)
				if err != nil {
					slog.Error("Failed to update user chat", "error", err.Error())
					return err
				}
			}
		}
	}
	return nil
}
