package controllers

import (
	"context"

	"gvnotes/internal/dto"
	"gvnotes/internal/services"
)

type ImageController struct {
	svc services.ImageService
}

func NewImageController(svc services.ImageService) *ImageController {
	return &ImageController{svc: svc}
}

func (c *ImageController) ListImagesByNote(noteID string) ([]dto.ImageItem, error) {
	return c.svc.ListByNote(context.Background(), noteID)
}

func (c *ImageController) GetImage(id string) (dto.ImageItem, error) {
	return c.svc.Get(context.Background(), id)
}

// SaveImage saves a base64-decoded image for a note.
// mimeType: e.g. "image/png", data: raw file bytes.
func (c *ImageController) SaveImage(noteID string, mimeType string, data []byte) (dto.ImageItem, error) {
	return c.svc.Save(context.Background(), noteID, mimeType, data)
}

func (c *ImageController) DeleteImage(id string) error {
	return c.svc.Delete(context.Background(), id)
}

// GetImagePath returns the absolute file path for an image filename.
// The frontend can use this with Wails runtime to load local files.
func (c *ImageController) GetImagePath(filename string) (string, error) {
	return c.svc.FilePath(filename)
}
