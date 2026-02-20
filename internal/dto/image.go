package dto

import "time"

type ImageItem struct {
	ID        string    `json:"id"`
	NoteID    string    `json:"noteId"`
	Filename  string    `json:"filename"`
	MimeType  string    `json:"mimeType"`
	CreatedAt time.Time `json:"createdAt"`
}
