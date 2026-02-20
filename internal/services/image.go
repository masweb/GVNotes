package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	db "gvnotes/db/generated"
	"gvnotes/internal/dto"
	apperrors "gvnotes/internal/errors"
)

type ImageService interface {
	ListByNote(ctx context.Context, noteID string) ([]dto.ImageItem, error)
	Get(ctx context.Context, id string) (dto.ImageItem, error)
	// Save stores the file on disk and registers it in the database.
	// data is the raw file bytes, mimeType is e.g. "image/png".
	Save(ctx context.Context, noteID string, mimeType string, data []byte) (dto.ImageItem, error)
	// Delete removes the database record and the file from disk.
	Delete(ctx context.Context, id string) error
	// FilePath returns the absolute path of an image file on disk.
	FilePath(filename string) (string, error)
}

type imageService struct {
	q         db.Querier
	imagesDir string
}

func NewImageService(q db.Querier, imagesDir string) ImageService {
	return &imageService{q: q, imagesDir: imagesDir}
}

func (s *imageService) ListByNote(ctx context.Context, noteID string) ([]dto.ImageItem, error) {
	rows, err := s.q.ListImagesByNote(ctx, noteID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.ImageItem, len(rows))
	for i, r := range rows {
		result[i] = toImageItem(r)
	}
	return result, nil
}

func (s *imageService) Get(ctx context.Context, id string) (dto.ImageItem, error) {
	img, err := s.q.GetImage(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.ImageItem{}, fmt.Errorf("%w: image %s", apperrors.ErrNotFound, id)
		}
		return dto.ImageItem{}, err
	}
	return toImageItem(img), nil
}

func (s *imageService) Save(ctx context.Context, noteID string, mimeType string, data []byte) (dto.ImageItem, error) {
	ext := extensionFromMIME(mimeType)
	filename := uuid.NewString() + ext
	fullPath := filepath.Join(s.imagesDir, filename)

	if err := os.WriteFile(fullPath, data, 0o644); err != nil {
		return dto.ImageItem{}, fmt.Errorf("write image file: %w", err)
	}

	img, err := s.q.CreateImage(ctx, db.CreateImageParams{
		ID:       uuid.NewString(),
		NoteID:   noteID,
		Filename: filename,
		MimeType: mimeType,
	})
	if err != nil {
		// Best-effort cleanup of the file if the DB insert fails.
		os.Remove(fullPath)
		return dto.ImageItem{}, err
	}

	return toImageItem(img), nil
}

func (s *imageService) Delete(ctx context.Context, id string) error {
	img, err := s.q.GetImage(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: image %s", apperrors.ErrNotFound, id)
		}
		return err
	}

	if err := s.q.DeleteImage(ctx, id); err != nil {
		return err
	}

	// Best-effort file removal; do not fail if the file is already gone.
	os.Remove(filepath.Join(s.imagesDir, img.Filename))
	return nil
}

func (s *imageService) FilePath(filename string) (string, error) {
	path := filepath.Join(s.imagesDir, filepath.Base(filename))
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("%w: image file %s", apperrors.ErrNotFound, filename)
	}
	return path, nil
}

func toImageItem(img db.Image) dto.ImageItem {
	return dto.ImageItem{
		ID:        img.ID,
		NoteID:    img.NoteID,
		Filename:  img.Filename,
		MimeType:  img.MimeType,
		CreatedAt: img.CreatedAt,
	}
}

func extensionFromMIME(mimeType string) string {
	switch strings.ToLower(mimeType) {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/svg+xml":
		return ".svg"
	default:
		return ".bin"
	}
}
