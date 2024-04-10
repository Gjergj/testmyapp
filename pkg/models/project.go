package models

const (
	MaxUploadFiles     = 35
	MaxFileSizeLimit   = 1 << 20 // 1 MB limit for one file
	MaxFileNameLength  = 255
	MaxUploadSize      = 1 << 20 // 1 MB limit for the entire request
	MaxProjectsPerUser = 5
)

type Project struct {
	ID             string `json:"id,omitempty"`
	UserID         string `json:"user_id,omitempty"`
	ProjectName    string `json:"project_name,omitempty"`
	Domain         string `json:"domain,omitempty"`
	CreationDate   string `json:"creation_date,omitempty"`
	LastDeployment string `json:"last_deployment,omitempty"`
	Deleted        bool   `json:"deleted,omitempty"`
	DeletedAt      string `json:"deleted_at,omitempty"`
	URL            string `json:"url,omitempty"`
}

func AllowedFileType(fileType string) bool {
	// use map instead of switch
	switch fileType {
	case ".jpeg", ".jpg", ".png", ".gif", ".html", ".htm", ".css", ".js", ".ico", ".json", ".svg", ".ttf", ".tiff", ".bmp", ".asp", ".txt", ".midi", ".mp4":
		return true
	}
	return false
}
