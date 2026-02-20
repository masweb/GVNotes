package controllers

import (
	"context"

	"gvnotes/internal/dto"
	apperrors "gvnotes/internal/errors"
	"gvnotes/internal/services"

	"errors"
)

type AuthController struct {
	svc services.AuthService
}

func NewAuthController(svc services.AuthService) *AuthController {
	return &AuthController{svc: svc}
}

// GetAuthStatus returns whether a password has been configured.
func (c *AuthController) GetAuthStatus() dto.AuthStatusResponse {
	set, err := c.svc.IsPasswordSet(context.Background())
	if err != nil {
		return dto.AuthStatusResponse{IsPasswordSet: false}
	}
	return dto.AuthStatusResponse{IsPasswordSet: set}
}

// SetPassword sets the app password for the first time.
func (c *AuthController) SetPassword(password string) dto.AuthResponse {
	if err := c.svc.SetPassword(context.Background(), password); err != nil {
		return dto.AuthResponse{Success: false, Error: friendlyError(err)}
	}
	return dto.AuthResponse{Success: true}
}

// VerifyPassword checks the provided password.
func (c *AuthController) VerifyPassword(password string) dto.AuthResponse {
	if err := c.svc.VerifyPassword(context.Background(), password); err != nil {
		return dto.AuthResponse{Success: false, Error: friendlyError(err)}
	}
	return dto.AuthResponse{Success: true}
}

func friendlyError(err error) string {
	switch {
	case errors.Is(err, apperrors.ErrUnauthenticated):
		return "incorrect password"
	case errors.Is(err, apperrors.ErrConflict):
		return "password already set"
	case errors.Is(err, apperrors.ErrInvalidInput):
		return "password cannot be empty"
	case errors.Is(err, apperrors.ErrNotFound):
		return "no password configured"
	default:
		return "internal error"
	}
}
