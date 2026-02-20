package controllers

import (
	"context"
	"errors"

	"gvnotes/internal/dto"
	apperrors "gvnotes/internal/errors"
	"gvnotes/internal/services"
)

type NoteController struct {
	svc services.NoteService
}

func NewNoteController(svc services.NoteService) *NoteController {
	return &NoteController{svc: svc}
}

// ListNotes returns notes at the root (notebookID == "") or inside a notebook.
func (c *NoteController) ListNotes(notebookID string) ([]dto.NoteListItem, error) {
	var nid *string
	if notebookID != "" {
		nid = &notebookID
	}
	return c.svc.List(context.Background(), nid)
}

func (c *NoteController) GetNote(id string) (dto.NoteDetail, error) {
	return c.svc.Get(context.Background(), id)
}

func (c *NoteController) CreateNote(req dto.CreateNoteRequest) (dto.NoteDetail, error) {
	return c.svc.Create(context.Background(), req)
}

func (c *NoteController) UpdateNoteTitle(id string, req dto.UpdateNoteTitleRequest) (dto.NoteDetail, error) {
	return c.svc.UpdateTitle(context.Background(), id, req)
}

func (c *NoteController) UpdateNoteContent(id string, req dto.UpdateNoteContentRequest) (dto.NoteDetail, error) {
	return c.svc.UpdateContent(context.Background(), id, req)
}

func (c *NoteController) UpdateNotePosition(id string, req dto.UpdatePositionRequest) error {
	return c.svc.UpdatePosition(context.Background(), id, req)
}

func (c *NoteController) DeleteNote(id string) error {
	return c.svc.Delete(context.Background(), id)
}

func noteFriendlyError(err error) string {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		return "note not found"
	case errors.Is(err, apperrors.ErrInvalidInput):
		return "title cannot be empty"
	default:
		return "internal error"
	}
}
