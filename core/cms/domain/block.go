package domain

// Block defines a CMS block
type Block struct {
	ID           int    `json:"id"`
	Identifier   string `json:"identifier"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	CreationTime string `json:"creation_time"`
	UpdateTime   string `json:"update_time"`
	Active       bool   `json:"active"`
}
