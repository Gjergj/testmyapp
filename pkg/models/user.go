package models

type User struct {
	ID       string `db:"id"  json:"id,omitempty"`
	Username string `db:"username" json:"username,omitempty"`
	Password string `db:"password" json:"password,omitempty"`
	Active   bool   `db:"active" json:"active,omitempty"`
	Deleted  bool   `db:"deleted" json:"deleted,omitempty"`
}

type ApiResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
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
