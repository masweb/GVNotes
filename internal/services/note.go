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

type NoteService interface {
	List(ctx context.Context, notebookID *string) ([]dto.NoteListItem, error)
	Get(ctx context.Context, id string) (dto.NoteDetail, error)
	Create(ctx context.Context, req dto.CreateNoteRequest) (dto.NoteDetail, error)
	UpdateTitle(ctx context.Context, id string, req dto.UpdateNoteTitleRequest) (dto.NoteDetail, error)
	UpdateContent(ctx context.Context, id string, req dto.UpdateNoteContentRequest) (dto.NoteDetail, error)
	UpdatePosition(ctx context.Context, id string, req dto.UpdatePositionRequest) error
	Delete(ctx context.Context, id string) error
}

type noteService struct {
	q db.Querier
}

func NewNoteService(q db.Querier) NoteService {
	return &noteService{q: q}
}

func (s *noteService) List(ctx context.Context, notebookID *string) ([]dto.NoteListItem, error) {
	rows, err := s.q.ListNotes(ctx, notebookID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.NoteListItem, len(rows))
	for i, r := range rows {
		result[i] = dto.NoteListItem{
			ID:        r.ID,
			Title:     r.Title,
			Position:  r.Position,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		}
	}
	return result, nil
}

func (s *noteService) Get(ctx context.Context, id string) (dto.NoteDetail, error) {
	n, err := s.q.GetNote(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.NoteDetail{}, fmt.Errorf("%w: note %s", apperrors.ErrNotFound, id)
		}
		return dto.NoteDetail{}, err
	}
	return toNoteDetail(n), nil
}

func (s *noteService) Create(ctx context.Context, req dto.CreateNoteRequest) (dto.NoteDetail, error) {
	if req.Title == "" {
		return dto.NoteDetail{}, fmt.Errorf("%w: title cannot be empty", apperrors.ErrInvalidInput)
	}
	n, err := s.q.CreateNote(ctx, db.CreateNoteParams{
		ID:         uuid.NewString(),
		NotebookID: req.NotebookID,
		Title:      req.Title,
		Content:    "{}",
		Position:   0,
	})
	if err != nil {
		return dto.NoteDetail{}, err
	}
	return toNoteDetail(n), nil
}

func (s *noteService) UpdateTitle(ctx context.Context, id string, req dto.UpdateNoteTitleRequest) (dto.NoteDetail, error) {
	if req.Title == "" {
		return dto.NoteDetail{}, fmt.Errorf("%w: title cannot be empty", apperrors.ErrInvalidInput)
	}
	n, err := s.q.UpdateNoteTitle(ctx, db.UpdateNoteTitleParams{
		ID:    id,
		Title: req.Title,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.NoteDetail{}, fmt.Errorf("%w: note %s", apperrors.ErrNotFound, id)
		}
		return dto.NoteDetail{}, err
	}
	return toNoteDetail(n), nil
}

func (s *noteService) UpdateContent(ctx context.Context, id string, req dto.UpdateNoteContentRequest) (dto.NoteDetail, error) {
	n, err := s.q.UpdateNoteContent(ctx, db.UpdateNoteContentParams{
		ID:      id,
		Content: req.Content,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.NoteDetail{}, fmt.Errorf("%w: note %s", apperrors.ErrNotFound, id)
		}
		return dto.NoteDetail{}, err
	}
	return toNoteDetail(n), nil
}

func (s *noteService) UpdatePosition(ctx context.Context, id string, req dto.UpdatePositionRequest) error {
	return s.q.UpdateNotePosition(ctx, db.UpdateNotePositionParams{
		ID:       id,
		Position: req.Position,
	})
}

func (s *noteService) Delete(ctx context.Context, id string) error {
	return s.q.DeleteNote(ctx, id)
}

func toNoteDetail(n db.Note) dto.NoteDetail {
	return dto.NoteDetail{
		ID:         n.ID,
		NotebookID: n.NotebookID,
		Title:      n.Title,
		Content:    n.Content,
		Position:   n.Position,
		CreatedAt:  n.CreatedAt,
		UpdatedAt:  n.UpdatedAt,
	}
}
