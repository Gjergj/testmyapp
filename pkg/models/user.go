package models

type User struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Active   bool   `json:"active,omitempty"`
	Deleted  bool   `json:"deleted,omitempty"`
}

type ApiResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type SignupResponse struct {
	ApiResponse
	UserID string `json:"user_id"`
}

type LoginResponse struct {
	ApiResponse
	UserID       string `json:"user_id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type CreateUserResponse struct {
	ApiResponse
	UserID string `json:"user_id"`
}

type GetProjectsResponse struct {
	ApiResponse
	Projects []*Project `json:"projects"`
}

type CreateProjectResponse struct {
	ApiResponse
	Project
}

type UploadFilesResponse struct {
	ApiResponse
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
