package controllers

import (
	"context"
	"errors"

	"gvnotes/internal/dto"
	apperrors "gvnotes/internal/errors"
	"gvnotes/internal/services"
)

type NotebookController struct {
	svc services.NotebookService
}

func NewNotebookController(svc services.NotebookService) *NotebookController {
	return &NotebookController{svc: svc}
}

// ListNotebooks returns notebooks at the root (parentID == "") or inside a parent.
func (c *NotebookController) ListNotebooks(parentID string) ([]dto.NotebookListItem, error) {
	var pid *string
	if parentID != "" {
		pid = &parentID
	}
	return c.svc.List(context.Background(), pid)
}

func (c *NotebookController) GetNotebook(id string) (dto.NotebookDetail, error) {
	return c.svc.Get(context.Background(), id)
}

func (c *NotebookController) CreateNotebook(req dto.CreateNotebookRequest) (dto.NotebookDetail, error) {
	return c.svc.Create(context.Background(), req)
}

func (c *NotebookController) UpdateNotebookTitle(id string, req dto.UpdateNotebookTitleRequest) (dto.NotebookDetail, error) {
	return c.svc.UpdateTitle(context.Background(), id, req)
}

func (c *NotebookController) UpdateNotebookPosition(id string, req dto.UpdatePositionRequest) error {
	return c.svc.UpdatePosition(context.Background(), id, req)
}

func (c *NotebookController) DeleteNotebook(id string) error {
	return c.svc.Delete(context.Background(), id)
}

func (c *NotebookController) MoveNotebook(id string, req dto.MoveNotebookRequest) (dto.NotebookDetail, error) {
	return c.svc.Move(context.Background(), id, req)
}

func notebookFriendlyError(err error) string {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		return "notebook not found"
	case errors.Is(err, apperrors.ErrInvalidInput):
		return "title cannot be empty"
	default:
		return "internal error"
	}
}
