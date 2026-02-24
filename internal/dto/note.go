package dto

import "time"

// NoteListItem is returned in list responses (no content).
type NoteListItem struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Position  int64     `json:"position"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NoteDetail is returned by GetNote (includes content).
type NoteDetail struct {
	ID         string    `json:"id"`
	NotebookID *string   `json:"notebookId"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Position   int64     `json:"position"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type CreateNoteRequest struct {
	NotebookID *string `json:"notebookId"`
	Title      string  `json:"title"`
}

type UpdateNoteTitleRequest struct {
	Title string `json:"title"`
}

type UpdateNoteContentRequest struct {
	Content string `json:"content"`
}

type MoveNoteRequest struct {
	NotebookID *string `json:"notebookId"`
}
