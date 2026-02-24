package dto

import "time"

// NotebookListItem is returned in list responses (no content, just metadata).
type NotebookListItem struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Position  int64     `json:"position"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NotebookDetail is returned by GetNotebook.
type NotebookDetail struct {
	ID        string    `json:"id"`
	ParentID  *string   `json:"parentId"`
	Title     string    `json:"title"`
	Position  int64     `json:"position"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateNotebookRequest struct {
	ParentID *string `json:"parentId"`
	Title    string  `json:"title"`
}

type UpdateNotebookTitleRequest struct {
	Title string `json:"title"`
}

type UpdatePositionRequest struct {
	Position int64 `json:"position"`
}

type MoveNotebookRequest struct {
	ParentID *string `json:"parentId"`
}
