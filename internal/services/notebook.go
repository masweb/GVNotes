package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	db "gvnotes/db/generated"
	"gvnotes/internal/dto"
	apperrors "gvnotes/internal/errors"
)

type NotebookService interface {
	List(ctx context.Context, parentID *string) ([]dto.NotebookListItem, error)
	Get(ctx context.Context, id string) (dto.NotebookDetail, error)
	Create(ctx context.Context, req dto.CreateNotebookRequest) (dto.NotebookDetail, error)
	UpdateTitle(ctx context.Context, id string, req dto.UpdateNotebookTitleRequest) (dto.NotebookDetail, error)
	UpdatePosition(ctx context.Context, id string, req dto.UpdatePositionRequest) error
	Delete(ctx context.Context, id string) error
	Move(ctx context.Context, id string, req dto.MoveNotebookRequest) (dto.NotebookDetail, error)
}

type notebookService struct {
	q db.Querier
}

func NewNotebookService(q db.Querier) NotebookService {
	return &notebookService{q: q}
}

func (s *notebookService) List(ctx context.Context, parentID *string) ([]dto.NotebookListItem, error) {
	rows, err := s.q.ListNotebooks(ctx, parentID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.NotebookListItem, len(rows))
	for i, r := range rows {
		result[i] = dto.NotebookListItem{
			ID:        r.ID,
			Title:     r.Title,
			Position:  r.Position,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		}
	}
	return result, nil
}

func (s *notebookService) Get(ctx context.Context, id string) (dto.NotebookDetail, error) {
	n, err := s.q.GetNotebook(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.NotebookDetail{}, fmt.Errorf("%w: notebook %s", apperrors.ErrNotFound, id)
		}
		return dto.NotebookDetail{}, err
	}
	return toNotebookDetail(n), nil
}

func (s *notebookService) Create(ctx context.Context, req dto.CreateNotebookRequest) (dto.NotebookDetail, error) {
	if req.Title == "" {
		return dto.NotebookDetail{}, fmt.Errorf("%w: title cannot be empty", apperrors.ErrInvalidInput)
	}

	existing, err := s.q.ListNotebooks(ctx, req.ParentID)
	if err != nil {
		return dto.NotebookDetail{}, err
	}
	var maxPos int64
	for _, nb := range existing {
		if nb.Position > maxPos {
			maxPos = nb.Position
		}
	}
	n, err := s.q.CreateNotebook(ctx, db.CreateNotebookParams{
		ID:       uuid.NewString(),
		ParentID: req.ParentID,
		Title:    req.Title,
		Position: maxPos + 1000,
	})
	if err != nil {
		return dto.NotebookDetail{}, err
	}
	return toNotebookDetail(n), nil
}

func (s *notebookService) UpdateTitle(ctx context.Context, id string, req dto.UpdateNotebookTitleRequest) (dto.NotebookDetail, error) {
	if req.Title == "" {
		return dto.NotebookDetail{}, fmt.Errorf("%w: title cannot be empty", apperrors.ErrInvalidInput)
	}
	n, err := s.q.UpdateNotebookTitle(ctx, db.UpdateNotebookTitleParams{
		ID:    id,
		Title: req.Title,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.NotebookDetail{}, fmt.Errorf("%w: notebook %s", apperrors.ErrNotFound, id)
		}
		return dto.NotebookDetail{}, err
	}
	return toNotebookDetail(n), nil
}

func (s *notebookService) UpdatePosition(ctx context.Context, id string, req dto.UpdatePositionRequest) error {
	return s.q.UpdateNotebookPosition(ctx, db.UpdateNotebookPositionParams{
		ID:       id,
		Position: req.Position,
	})
}

func (s *notebookService) Delete(ctx context.Context, id string) error {
	return s.q.DeleteNotebook(ctx, id)
}

func (s *notebookService) Move(ctx context.Context, id string, req dto.MoveNotebookRequest) (dto.NotebookDetail, error) {
	if req.ParentID != nil && *req.ParentID == id {
		return dto.NotebookDetail{}, fmt.Errorf("%w: cannot move notebook into itself", apperrors.ErrInvalidInput)
	}
	existing, err := s.q.ListNotebooks(ctx, req.ParentID)
	if err != nil {
		return dto.NotebookDetail{}, err
	}
	var maxPos int64
	for _, nb := range existing {
		if nb.Position > maxPos {
			maxPos = nb.Position
		}
	}
	n, err := s.q.MoveNotebook(ctx, db.MoveNotebookParams{
		ParentID: req.ParentID,
		Position: maxPos + 1000,
		ID:       id,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.NotebookDetail{}, fmt.Errorf("%w: notebook %s", apperrors.ErrNotFound, id)
		}
		return dto.NotebookDetail{}, err
	}
	return toNotebookDetail(n), nil
}

func toNotebookDetail(n db.Notebook) dto.NotebookDetail {
	return dto.NotebookDetail{
		ID:        n.ID,
		ParentID:  n.ParentID,
		Title:     n.Title,
		Position:  n.Position,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}
