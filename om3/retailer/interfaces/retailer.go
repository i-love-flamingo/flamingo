package interfaces

type (
	// RetailerService is used to retrieve retailer
	RetailerService interface {
		ByID() Retailer
	}

	// Retailer defines the retailer model
	Retailer interface {
		ID() string
		Name() string
	}
)
