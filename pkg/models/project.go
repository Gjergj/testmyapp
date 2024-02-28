package models

type Project struct {
	ID             string `db:"id" json:"id,omitempty"`
	UserID         string `db:"user_id" json:"user_id,omitempty"`
	ProjectName    string `db:"project_name" json:"project_name,omitempty"`
	Domain         string `db:"domain" json:"domain,omitempty"`
	CreationDate   string `db:"creation_date" json:"creation_date,omitempty"`
	LastDeployment string `db:"last_deployment" json:"last_deployment,omitempty"`
	Deleted        bool   `db:"deleted" json:"deleted,omitempty"`
	DeletedAt      string `db:"deleted_at" json:"deleted_at,omitempty"`
	URL            string `json:"url,omitempty"`
}
