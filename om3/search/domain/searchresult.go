package domain

import (
	"flamingo/core/product/domain"
)

type (
	SearchResult struct {
		Results struct {
			Retailer struct {
				MetaData MetaData      `json:"metaData"`
				Facets   []interface{} `json:"facets"`
				Filters  []interface{} `json:"filters"`
				PageInfo PageInfo      `json:"pageInfo"`
				Hits     []interface{} `json:"hits"`
			} `json:"retailer"`
			Location struct {
				MetaData MetaData      `json:"metaData"`
				Facets   []interface{} `json:"facets"`
				Filters  []interface{} `json:"filters"`
				PageInfo PageInfo      `json:"pageInfo"`
				Hits     []interface{} `json:"hits"`
			} `json:"location"`
			Brand struct {
				MetaData MetaData      `json:"metaData"`
				Facets   []interface{} `json:"facets"`
				Filters  []interface{} `json:"filters"`
				PageInfo PageInfo      `json:"pageInfo"`
				Hits     []interface{} `json:"hits"`
			} `json:"brand"`
			Product struct {
				MetaData MetaData      `json:"metaData"`
				Facets   []interface{} `json:"facets"`
				Filters  []interface{} `json:"filters"`
				PageInfo PageInfo      `json:"pageInfo"`
				Hits     []struct {
					Document   domain.Product `json:"document"`
					Highlights struct {
					} `json:"highlights"`
				} `json:"hits"`
			} `json:"product"`
		} `json:"results"`
	}

	MetaData struct {
		TotalHits    int    `json:"totalHits"`
		Took         int    `json:"took"`
		CurrentQuery string `json:"currentQuery"`
		FacetMapping []struct {
			DocumentType string        `json:"documentType"`
			FacetNames   []interface{} `json:"facetNames"`
		} `json:"facetMapping"`
		SortMapping []struct {
			DocumentType string `json:"documentType"`
			Sorts        struct {
			} `json:"sorts"`
		} `json:"sortMapping"`
	}

	PageInfo struct {
		CurrentPage      int `json:"currentPage"`
		TotalPages       int `json:"totalPages"`
		VisiblePageLinks int `json:"visiblePageLinks"`
		Padding          int `json:"padding"`
	}
)
