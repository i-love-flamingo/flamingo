package domain

// Page defines a CMS Page
type Page struct {
	ID              int    `json:"id"`
	Identifier      string `json:"identifier"`
	Title           string `json:"title"`
	PageLayout      string `json:"page_layout"`
	MetaTitle       string `json:"meta_title"`
	MetaKeywords    string `json:"meta_keywords"`
	MetaDescription string `json:"meta_description"`
	ContentHeading  string `json:"content_heading"`
	Content         string `json:"content"`
	CreationTime    string `json:"creation_time"`
	UpdateTime      string `json:"update_time"`
	SortOrder       string `json:"sort_order"`
	LayoutUpdateXML string `json:"layout_update_xml"`
	Active          bool   `json:"active"`
}
