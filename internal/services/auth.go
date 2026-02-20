package services

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"

	db "gvnotes/db/generated"
	apperrors "gvnotes/internal/errors"

	"golang.org/x/crypto/argon2"
)

const passwordKey = "auth.password_hash"

// argon2id parameters (OWASP recommended minimums).
const (
	argonTime    = 1
	argonMemory  = 64 * 1024 // 64 MB
	argonThreads = 4
	argonKeyLen  = 32
	saltLen      = 16
)

//go:generate mockery --name AuthService
type AuthService interface {
	IsPasswordSet(ctx context.Context) (bool, error)
	SetPassword(ctx context.Context, password string) error
	VerifyPassword(ctx context.Context, password string) error
}

type authService struct {
	q db.Querier
}

func NewAuthService(q db.Querier) AuthService {
	return &authService{q: q}
}

func (s *authService) IsPasswordSet(ctx context.Context) (bool, error) {
	_, err := s.q.GetSetting(ctx, passwordKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *authService) SetPassword(ctx context.Context, password string) error {
	if password == "" {
		return fmt.Errorf("%w: password cannot be empty", apperrors.ErrInvalidInput)
	}

	set, err := s.IsPasswordSet(ctx)
	if err != nil {
		return err
	}
	if set {
		return fmt.Errorf("%w: password is already set", apperrors.ErrConflict)
	}

	hash, err := hashPassword(password)
	if err != nil {
		return err
	}

	return s.q.UpsertSetting(ctx, db.UpsertSettingParams{
		Key:   passwordKey,
		Value: hash,
	})
}

func (s *authService) VerifyPassword(ctx context.Context, password string) error {
	encoded, err := s.q.GetSetting(ctx, passwordKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: no password set", apperrors.ErrNotFound)
		}
		return err
	}

	if !checkPassword(password, encoded) {
		return apperrors.ErrUnauthenticated
	}
	return nil
}

// hashPassword returns a base64-encoded "salt:hash" string.
func hashPassword(password string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	encoded := base64.StdEncoding.EncodeToString(salt) + ":" + base64.StdEncoding.EncodeToString(hash)
	return encoded, nil
}

// checkPassword verifies a password against a stored "salt:hash" string.
func checkPassword(password, encoded string) bool {
	saltB64, hashB64, ok := splitEncoded(encoded)
	if !ok {
		return false
	}

	salt, err := base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		return false
	}
	expectedHash, err := base64.StdEncoding.DecodeString(hashB64)
	if err != nil {
		return false
	}

	actualHash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	return subtle.ConstantTimeCompare(actualHash, expectedHash) == 1
}

func splitEncoded(encoded string) (salt, hash string, ok bool) {
	for i, c := range encoded {
		if c == ':' {
			return encoded[:i], encoded[i+1:], true
		}
	}
	return "", "", false
}
