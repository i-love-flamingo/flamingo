package interfaces

type PageService interface {
	Get(string) Page
}

type Page interface {
	Name() string
	Content() string
}
