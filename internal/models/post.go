// Package models
package models

import "time"

// Post represents a blog post in our domain. We include some JSON-typed fields
// which will be stored in MySQL JSON columns.
type Post struct {
	ID        int64          `json:"id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Tags      []string       `json:"tags"`
	Metadata  map[string]any `json:"metadata"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}
