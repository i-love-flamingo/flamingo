package interfaces

// PageService defines the page-getter service
type PageService interface {
	Get(string) Page
}

// Page defines what a CMS page object looks like
type Page interface {
	Name() string
	Content() string
}
