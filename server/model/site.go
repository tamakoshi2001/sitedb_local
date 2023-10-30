package model

import (
	"time"
)

type (
	// A Site represents a Site which has url, title and summary.
	Site struct {
		URL       string    `json:"url"`
		Title     string    `json:"title"`
		Summary   string    `json:"summary"`
		CreatedAt time.Time `json:"created_at"`
	}

	Vector struct {
		URL    string     `json:"url"`
		Vector *[]float64 `json:"vector"`
	}

	// A CreateSiteRequest represents a request for creating a Site.
	CreateSiteRequest struct {
		URL string `json:"url"`
	}

	// A CreateSiteResponse represents a response for creating a Site.
	CreateSiteResponse struct {
		Site Site `json:"site"`
	}

	// A ReadSiteRequest represents a request for reading Sites.
	ReadSiteRequest struct {
		Query string `json:"query"`
		// Size   int64 `json:"size"`
		// How   string `json:"how"`
	}

	// A ReadSiteResponse represents a response for reading Sites.
	ReadSiteResponse struct {
		SiteS []Site `json:"Sites"`
	}

	// A DeleteSiteRequest represents a request for deleting a Site.
	DeleteSiteRequest struct {
		URL string `json:"url"`
	}

	// A DeleteSiteResponse represents a response for deleting a Site.
	DeleteSiteResponse struct {
		Response string `json:"response"`
	}
)
