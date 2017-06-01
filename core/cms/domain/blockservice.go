package domain

// PageService defines the page-getter service
type BlockService interface {
	Get(string) (*Block, error)
}
