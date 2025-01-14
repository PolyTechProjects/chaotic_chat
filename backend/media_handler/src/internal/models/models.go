package models

import (
	"github.com/google/uuid"
)

type Media struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	ObjectType string    `gorm:"not null;check:object_type <> ''"`
	ObjectId   string    `gorm:"not null;check:object_id <> ''"`
	FileId     string    `gorm:"not null;check:file_id <> ''"`
}

func New(id uuid.UUID, objectType string, objectId string, fileId string) *Media {
	return &Media{ID: id, ObjectType: objectType, ObjectId: objectId, FileId: fileId}
}

type SeaweedFSAssignResponse struct {
	Count     int    `json:"count"`
	Fid       string `json:"fid"`
	Url       string `json:"url"`
	PublicUrl string `json:"publicUrl"`
}

type SeaweedFSLookupResponse struct {
	VolumeId  int             `json:"volumeId"`
	Locations []PublicUrlXUrl `json:"locations"`
}

type PublicUrlXUrl struct {
	PublicUrl string `json:"publicUrl"`
	Url       string `json:"url"`
}

type UploadMediaRequest struct {
	MessageId string `json:"messageId"`
}

type MessageIdXFileId struct {
	MessageId uuid.UUID `json:"messageId"`
	FileId    uuid.UUID `json:"fileId"`
}

type UploadMediaResponse struct {
	ObjectType string `json:"objectType"`
	ObjectId   string `json:"objectId"`
	FileId     string `json:"fileId"`
}
