package interfaces

type (
	BrandService interface {
		ByID(string) Brand
	}

	Brand interface {
		ID() string
		Name() string
	}
)
