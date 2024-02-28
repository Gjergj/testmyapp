package models

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
