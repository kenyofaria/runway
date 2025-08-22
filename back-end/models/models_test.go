package models

import (
	"testing"
)

func TestReview_ToReviewResponse(t *testing.T) {
	// Test case for a successful conversion
	t.Run("successful conversion", func(t *testing.T) {
		review := Review{
			ID: struct {
				Label string `json:"label"`
			}(struct{ Label string }{Label: "12345"}),
			Author: struct {
				Name struct {
					Label string `json:"label"`
				} `json:"name"`
			}(struct{ Name struct{ Label string } }{Name: struct{ Label string }{Label: "John Doe"}}),
			Content: struct {
				Label string `json:"label"`
			}(struct{ Label string }{Label: "This is a great app."}),
			Rating: struct {
				Label string `json:"label"`
			}(struct{ Label string }{Label: "5"}),
			Timestamp: struct {
				Label string `json:"label"`
			}(struct{ Label string }{Label: "2025-08-21T10:00:00Z"}),
		}

		resp, err := review.ToReviewResponse()
		if err != nil {
			t.Fatalf("ToReviewResponse() returned an unexpected error: %v", err)
		}

		if resp.ID != "12345" {
			t.Errorf("Expected ID '12345', but got '%s'", resp.ID)
		}
		if resp.Author != "John Doe" {
			t.Errorf("Expected Author 'John Doe', but got '%s'", resp.Author)
		}
		if resp.Content != "This is a great app." {
			t.Errorf("Expected Content 'This is a great app.', but got '%s'", resp.Content)
		}
		if resp.Score != 5 {
			t.Errorf("Expected Score 5, but got %d", resp.Score)
		}
		if resp.Time != "2025-08-21T10:00:00Z" {
			t.Errorf("Expected Time '2025-08-21T10:00:00Z', but got '%s'", resp.Time)
		}
	})

	// Test case for a failed conversion (e.g., invalid score)
	t.Run("failed conversion with invalid score", func(t *testing.T) {
		review := Review{
			ID: struct {
				Label string `json:"label"`
			}(struct{ Label string }{Label: "67890"}),
			Author: struct {
				Name struct {
					Label string `json:"label"`
				} `json:"name"`
			}(struct{ Name struct{ Label string } }{Name: struct{ Label string }{Label: "Jane Doe"}}),
			Content: struct {
				Label string `json:"label"`
			}(struct{ Label string }{Label: "Bad rating score."}),
			Rating: struct {
				Label string `json:"label"`
			}(struct{ Label string }{Label: "five"}), // Invalid score
			Timestamp: struct {
				Label string `json:"label"`
			}(struct{ Label string }{Label: "2025-08-21T11:00:00Z"}),
		}

		_, err := review.ToReviewResponse()
		if err == nil {
			t.Fatal("ToReviewResponse() was expected to return an error, but it did not.")
		}
	})
}
