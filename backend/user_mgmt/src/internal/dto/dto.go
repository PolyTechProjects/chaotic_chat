package dto

type UpdateInfoRequest struct {
	UserId      string `json:"user_id"`
	Name        string `json:"name"`
	UrlTag      string `json:"url_tag"`
	Description string `json:"description"`
}

type UpdateInfoResponse struct {
	UserId string `json:"user_id"`
}

type UploadProfilePicRequest struct {
	UserId string `json:"user_id"`
	FileId string `json:"file_id"`
}

type GetUserResponse struct {
	UserId      string `json:"user_id"`
	Name        string `json:"name"`
	UrlTag      string `json:"url_tag"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}

type DeleteUserRequest struct {
	UserId string `json:"user_id"`
}
