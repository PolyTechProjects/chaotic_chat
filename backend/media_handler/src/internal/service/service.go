package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/config"
	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/internal/models"
	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/internal/repository"
	"github.com/google/uuid"
)

type MediaHandlerService struct {
	mediaHandlerRepository *repository.MediaHandlerRepository
	masterUrl              string
}

func New(mediaHandlerRepository *repository.MediaHandlerRepository, cfg *config.Config) *MediaHandlerService {
	return &MediaHandlerService{
		mediaHandlerRepository: mediaHandlerRepository,
		masterUrl:              fmt.Sprintf("%s:%d", cfg.SeaweedFS.MasterIp, cfg.SeaweedFS.MasterPort),
	}
}

func (m *MediaHandlerService) UploadMedia(
	objectType string,
	objectId uuid.UUID,
	file multipart.File,
	fileHeader *multipart.FileHeader) (uuid.UUID, error) {
	fileId, err := m.assignFileToSeaweedFS(file, fileHeader.Filename)
	if err != nil {
		slog.Error(err.Error())
		return uuid.Nil, err
	}
	id := uuid.New()
	media := models.New(id, objectType, objectId.String(), fileId)
	err = m.mediaHandlerRepository.Save(media)
	if err != nil {
		slog.Error(err.Error())
		return uuid.Nil, err
	}
	if fileHeader.Size > 4194304 {
		err = fmt.Errorf("fileHeader.Size: %d is too large (over 4Mb)", fileHeader.Size)
		return uuid.Nil, err
	}

	switch objectType {
	case "OBJECT_MESSAGE":
		messageId := objectId
		mf := models.MessageIdXFileId{
			MessageId: messageId,
			FileId:    media.ID,
		}
		bytes, err := json.Marshal(mf)
		if err != nil {
			slog.Error(err.Error())
			return uuid.Nil, err
		}
		err = m.mediaHandlerRepository.PublishInFileLoadedChannel(bytes)
		if err != nil {
			slog.Error(err.Error())
			return uuid.Nil, err
		}
		return uuid.Nil, nil
	case "OBJECT_USER":
		return media.ID, nil
	case "OBJECT_CHAT":
		return media.ID, nil
	}
	return uuid.Nil, nil
}

func (m *MediaHandlerService) GetMedia(id uuid.UUID) ([]byte, error) {
	fileId, volumeAddress, err := m.lookUpForFileIdAndVolumeAddress(id)
	if err != nil {
		return nil, err
	}
	slog.Info(fileId)
	slog.Info(volumeAddress)
	res, err := http.Get(fmt.Sprintf("http://%s/%s", volumeAddress, fileId))
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (m *MediaHandlerService) DeleteMedia(id uuid.UUID) error {
	fileId, volumeAddress, err := m.lookUpForFileIdAndVolumeAddress(id)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s/%s", volumeAddress, fileId), nil)
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	err = m.mediaHandlerRepository.DeleteById(id)
	if err != nil {
		return err
	}
	return nil
}

func (m *MediaHandlerService) lookUpForFileIdAndVolumeAddress(id uuid.UUID) (string, string, error) {
	media, err := m.mediaHandlerRepository.FindById(id)
	if err != nil {
		return "", "", err
	}
	volumeId := strings.Split(media.FileId, ",")[0]

	url, err := m.mediaHandlerRepository.GetVolumeIp(volumeId)
	if err == nil {
		return media.FileId, url, nil
	}

	lookupResponse := &models.SeaweedFSLookupResponse{}
	res, err := http.Get(fmt.Sprintf("http://%s/dir/lookup?volumeId=%s", m.masterUrl, volumeId))
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(lookupResponse)

	url = lookupResponse.Locations[0].Url
	err = m.mediaHandlerRepository.CacheVolumeIp(volumeId, url)
	if err != nil {
		return "", "", err
	}

	return media.FileId, url, nil
}

func (m *MediaHandlerService) assignFileToSeaweedFS(file io.Reader, fileName string) (string, error) {
	assignResponse := &models.SeaweedFSAssignResponse{}
	res, err := http.Get(fmt.Sprintf("http://%s/dir/assign", m.masterUrl))
	if err != nil {
		return "", err
	}
	json.NewDecoder(res.Body).Decode(assignResponse)
	defer res.Body.Close()

	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)
	form, err := w.CreateFormFile("file", fileName)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(form, file)
	if err != nil {
		return "", err
	}
	w.Close()

	addr := fmt.Sprintf("http://%s/%s", assignResponse.Url, assignResponse.Fid)
	slog.Info(fmt.Sprintf("File URL: %v", addr))
	req, err := http.NewRequest("POST", addr, b)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	err = m.mediaHandlerRepository.CacheVolumeIp(strings.Split(assignResponse.Fid, ",")[0], assignResponse.Url)
	if err != nil {
		return "", err
	}

	return assignResponse.Fid, nil
}
