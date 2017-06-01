package domain

type Page struct {
	ID                    int    `json:"id"`
	Identifier            string `json:"identifier"`
	Title                 string `json:"title"`
	PageLayout            string `json:"pageLayout"`
	MetaKeywords          string `json:"metaKeywords"`
	MetaDescription       string `json:"metaDescription"`
	ContentHeading        string `json:"contentHeading"`
	Content               string `json:"content"`
	CreationTime          string `json:"creationTime"`
	UpdateTime            string `json:"updateTime"`
	SortOrder             string `json:"sortOrder"`
	LayoutUpdateXML       string `json:"layoutUpdateXml"`
	CustomTheme           string `json:"customTheme"`
	CustomRootTemplate    string `json:"customRootTemplate"`
	CustomLayoutUpdateXML string `json:"customLayoutUpdateXml"`
	CustomThemeFrom       string `json:"customThemeFrom"`
	CustomThemeTo         string `json:"customThemeTo"`
	IsActive              bool   `json:"isActive"`
}
