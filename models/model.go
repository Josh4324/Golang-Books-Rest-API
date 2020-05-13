package models

// User schema of the user table
type Book struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	PublishedAt string `json:"published_at"`
}
