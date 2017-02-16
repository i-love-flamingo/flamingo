package interfaces

type (
	RetailerService interface {
		ByID() Retailer
	}

	Retailer interface {
		ID() string
		Name() string
	}
)
