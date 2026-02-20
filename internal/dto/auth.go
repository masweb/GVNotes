package dto

// AuthStatusResponse tells the frontend whether a password has been set.
type AuthStatusResponse struct {
	IsPasswordSet bool `json:"isPasswordSet"`
}

// AuthResponse is returned after SetPassword or VerifyPassword.
type AuthResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}
