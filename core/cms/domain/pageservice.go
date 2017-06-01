package domain

type PageService interface {
	Get(string) (*Page, error)
}
