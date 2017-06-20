package domain

// Block defines a CMS block
type Block struct {
	ID           int    `json:"id"`
	Identifier   string `json:"identifier"`
	Title        string `json:"title,omitempty"`
	Content      string `json:"content,omitempty"`
	CreationTime string `json:"creation_time,omitempty"`
	UpdateTime   string `json:"update_time,omitempty"`
	Active       bool   `json:"active,omitempty"`
}
