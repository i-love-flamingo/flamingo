package domain

type (
	// Brand defines a product brand
	Brand struct {
		ID               string             `json:"id"`
		SegmentRelevance map[string]float64 `json:"segmentRelevance"`
		Logo             BrandImage         `json:"logo"`
		HeroImage        BrandImage         `json:"heroImage"`
		AdditionalMedia  []interface{}      `json:"additionalMedia"`
		WithDetailPage   bool               `json:"withDetailPage"`
		DefaultRelevance int                `json:"defaultRelevance"`
		LastUpdate       string             `json:"lastUpdate"`
		Title            string             `json:"title"`
		ShortTitle       string             `json:"shortTitle"`
		Teaser           string             `json:"teaser"`
		Introduction     string             `json:"introduction"`
		Description      string             `json:"description"`
		Keywords         string             `json:"keywords"`
		LocaleRelevance  int                `json:"localeRelevance"`
		Language         string             `json:"language"`
		Url              string             `json:"url"`
	}

	// BrandImage can be displayed
	BrandImage struct {
		ID    int    `json:"id"`
		Hash  string `json:"hash"`
		Url   string `json:"url"`
		Fixed string `json:"fixed"`
		One1  string `json:"1-1"`
		One69 string `json:"16-9"`
	}
)
