package dtos

import "github.com/google/uuid"

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@mail.com"`
	Password string `json:"password" binding:"required" example:"Password123"`
}

type AuthResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Token  string    `json:"token"`
}
