package domain

// Page defines a CMS Page
type Page struct {
	ID              int              `json:"id"`
	Identifier      string           `json:"identifier"`
	Title           string           `json:"title"`
	PageLayout      string           `json:"page_layout,omitempty"`
	MetaTitle       string           `json:"meta_title,omitempty"`
	MetaKeywords    string           `json:"meta_keywords,omitempty"`
	MetaDescription string           `json:"meta_description,omitempty"`
	ContentHeading  string           `json:"content_heading,omitempty"`
	Content         string           `json:"content,omitempty"`
	CreationTime    string           `json:"creation_time,omitempty"`
	UpdateTime      string           `json:"update_time,omitempty"`
	SortOrder       string           `json:"sort_order,omitempty"`
	LayoutUpdateXML string           `json:"layout_update_xml,omitempty"`
	Active          bool             `json:"active,omitempty"`
	BluefootEnabled bool             `json:"bluefoot_enabled,omitempty"`
	BluefootContent []BluefootEntity `json:"bluefoot_content,omitempty"`
}

type BluefootEntity struct {
	Type     string            `json:"type"`
	Children []BluefootContent `json:"children"`
}

type BluefootContent struct {
	ContentType string                       `json:"contentType"`
	FormData    map[string]string            `json:"formData"`
	Children    map[string][]BluefootContent `json:"children"`
}
