package book

import "notes-go/internal/author"

type Book struct {
	Id      string          `json:"id"`
	Name    string          `json:"name"`
	Age     int             `json:"age"`
	Authors []author.Author `json:"authors"`
}
