package domain

// Block defines a CMS block
type Block struct {
	ID           int    `json:"id"`
	Identifier   string `json:"identifier"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	CreationTime string `json:"creationTime"`
	UpdateTime   string `json:"updateTime"`
	IsActive     bool   `json:"isActive"`
}
