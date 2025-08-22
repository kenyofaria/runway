package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// CustomTime is a custom type for unmarshaling and marshaling date strings.
type CustomTime struct {
	t time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	if s == "null" || s == "" {
		return nil
	}
	// The layout "2006-01-02" is the reference date format for Go's time package.
	ct.t, err = time.Parse("2006-01-02", s)
	if err != nil {
		// As a fallback, try parsing the full RFC3339 format, which is how Go saves time.
		ct.t, err = time.Parse(time.RFC3339, s)
	}
	return err
}

// MarshalJSON implements the json.Marshaler interface.
// It formats the time back into the "YYYY-MM-DD" string format.
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, ct.t.Format("2006-01-02"))), nil
}

// Time provides access to the underlying time.Time value.
func (ct CustomTime) Time() time.Time {
	return ct.t
}

// Root is the top-level struct for the entire JSON payload.
type Root struct {
	Feed struct {
		Entries []App `json:"entry"`
	} `json:"feed"`
}

// Entry represents a single application's data.
type App struct {
	IMName        LabelField  `json:"im:name"`
	IMImages      []Image     `json:"im:image"`
	Summary       LabelField  `json:"summary"`
	IMPrice       Price       `json:"im:price"`
	IMContentType ContentType `json:"im:contentType"`
	Rights        LabelField  `json:"rights"`
	Title         LabelField  `json:"title"`
	// Use json.RawMessage to handle the polymorphic 'link' field.
	LinkRaw       json.RawMessage `json:"link"`
	LinkSingle    Link            `json:"-"`
	LinkMulti     []Link          `json:"-"`
	ID            AppID           `json:"id"`
	IMArtist      Artist          `json:"im:artist"`
	Category      Category        `json:"category"`
	IMReleaseDate ReleaseDate     `json:"im:releaseDate"`
}

type AppResponse struct {
	ID          string `json:"id"`
	AppID       string `json:"app_id"`
	BundleID    string `json:"bundle_id"`
	Author      string `json:"author"`
	ReleaseDate string `json:"release_date"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	ArtworkURL  string `json:"artwork_url"`
	URL         string `json:"url"`
	Summary     string `json:"summary"`
	Price       string `json:"price"`
	Rights      string `json:"rights"`
	Title       string `json:"title"`
}

func (a *App) ToAppResponse() (*AppResponse, error) {
	var artworkURL string
	if len(a.IMImages) > 0 {
		artworkURL = a.IMImages[0].Label
	}

	var appURL string
	if len(a.LinkMulti) > 0 {
		for _, link := range a.LinkMulti {
			if link.Attributes.Rel == "alternate" && link.Attributes.Type == "text/html" {
				appURL = link.Attributes.Href
				break
			}
		}
	} else if a.LinkSingle.Attributes.Rel == "alternate" && a.LinkSingle.Attributes.Type == "text/html" {
		appURL = a.LinkSingle.Attributes.Href
	}

	return &AppResponse{
		ID:          a.ID.Attributes.ID,
		AppID:       a.ID.Attributes.ID,
		BundleID:    a.ID.Attributes.BundleID,
		Author:      a.IMArtist.Label,
		ReleaseDate: a.IMReleaseDate.Attributes.Label,
		Name:        a.IMName.Label,
		Category:    a.Category.Attributes.Label,
		ArtworkURL:  artworkURL,
		URL:         appURL,
		Summary:     a.Summary.Label,
		Price:       a.IMPrice.Label,
		Rights:      a.Rights.Label,
		Title:       a.Title.Label,
	}, nil
}

// LabelField is a generic struct for fields with just a "label" key.
type LabelField struct {
	Label string `json:"label"`
}

// Image represents a single image with its attributes.
type Image struct {
	Label      string `json:"label"`
	Attributes struct {
		Height string `json:"height"`
	} `json:"attributes"`
}

// Price represents the application's price.
type Price struct {
	Label      string `json:"label"`
	Attributes struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"attributes"`
}

// ContentType represents the content type of the application.
type ContentType struct {
	Attributes struct {
		Term  string `json:"term"`
		Label string `json:"label"`
	} `json:"attributes"`
}

// AppID represents the application's unique identifier.
type AppID struct {
	Label      string `json:"label"`
	Attributes struct {
		ID       string `json:"im:id"`
		BundleID string `json:"im:bundleId"`
	} `json:"attributes"`
}

// Artist represents the application's artist/developer.
type Artist struct {
	Label      string `json:"label"`
	Attributes struct {
		Href string `json:"href"`
	} `json:"attributes"`
}

// Category represents the application's category.
type Category struct {
	Attributes struct {
		ID     string `json:"im:id"`
		Term   string `json:"term"`
		Scheme string `json:"scheme"`
		Label  string `json:"label"`
	} `json:"attributes"`
}

// ReleaseDate represents the application's release date.
type ReleaseDate struct {
	Label      string `json:"label"`
	Attributes struct {
		Label string `json:"label"`
	} `json:"attributes"`
}

// Link represents a link related to the application.
type Link struct {
	Attributes struct {
		Rel         string `json:"rel"`
		Type        string `json:"type"`
		Href        string `json:"href"`
		Title       string `json:"title,omitempty"`
		IMDuration  string `json:"im:duration,omitempty"`
		IMAssetType string `json:"im:assetType,omitempty"`
	} `json:"attributes"`
}

// The UnmarshalJSON method for Entry handles the dynamic `link` field.
func (e *App) UnmarshalJSON(data []byte) error {
	// A temporary struct to avoid infinite recursion.
	type Alias App
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	// Unmarshal all fields except 'link' into the temporary struct.
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Try to unmarshal the raw 'link' message as a single Link struct.
	var single Link
	if err := json.Unmarshal(e.LinkRaw, &single); err == nil {
		e.LinkSingle = single
		return nil
	}

	// If it fails, assume it's an array and try to unmarshal it as such.
	var multi []Link
	if err := json.Unmarshal(e.LinkRaw, &multi); err == nil {
		e.LinkMulti = multi
		return nil
	}

	// If both attempts fail, return an error.
	return fmt.Errorf("failed to unmarshal 'link' field")
}

// APIResponse represents the complete top-level JSON structure for the app list.
//type APIResponse struct {
//	Feed Feed `json:"feed"`
//}

//// Feed represents the nested "feed" object for the app list.
//type Feed struct {
//	Results []App `json:"entry"`
//}
//
//// App represents a single app with its metadata.
//type App struct {
//	ID          string     `json:"id"`
//	Author      string     `json:"artistName"`
//	ReleaseDate CustomTime `json:"releaseDate"`
//	Name        string     `json:"name"`
//	Kind        string     `json:"kind"`
//	ArtworkURL  string     `json:"artworkUrl100"`
//	URL         string     `json:"url"`
//}

// ReviewFeed represents the top-level JSON structure for app reviews.
type ReviewFeed struct {
	Feed ReviewsFeed `json:"feed"`
}

// ReviewsFeed represents the "feed" object for reviews.
type ReviewsFeed struct {
	Entries []Review `json:"entry"`
}

// Review represents a single app review.
type Review struct {
	ID struct {
		Label string `json:"label"`
	} `json:"id"`
	Author struct {
		Name struct {
			Label string `json:"label"`
		} `json:"name"`
	} `json:"author"`
	Content struct {
		Label string `json:"label"`
	} `json:"content"`
	Rating struct {
		Label string `json:"label"`
	} `json:"im:rating"`
	Timestamp struct {
		Label string `json:"label"`
	} `json:"updated"`
}

// AppReviews stores a collection of reviews for a specific app.
type AppReviews struct {
	mu      sync.Mutex // Mutex to protect data access
	Reviews []Review   `json:"reviews"`
}

// ReviewResponse is the simplified struct used for the API's public response.
type ReviewResponse struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Author  string `json:"author"`
	Score   int    `json:"score"`
	Time    string `json:"time"`
}

// ToReviewResponse converts a Review struct to a simplified ReviewResponse struct.
func (r *Review) ToReviewResponse() (*ReviewResponse, error) {
	score, err := strconv.Atoi(r.Rating.Label)
	if err != nil {
		return nil, fmt.Errorf("failed to convert rating to integer: %w", err)
	}

	return &ReviewResponse{
		ID:      r.ID.Label,
		Content: r.Content.Label,
		Author:  r.Author.Name.Label,
		Score:   score,
		Time:    r.Timestamp.Label,
	}, nil
}
