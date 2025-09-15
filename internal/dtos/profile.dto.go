package dtos

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type UpdateProfileRequest struct {
	FirstName   *string               `form:"firstname" json:"firstname" example:"Farid"`
	LastName    *string               `form:"lastname" json:"lastname" example:"Darari"`
	PhoneNumber *string               `form:"phone_number" json:"phone_number" example:"08123456789"`
	Avatar      *multipart.FileHeader `form:"avatar"`
}

type ProfileResponse struct {
	UserID      uuid.UUID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FirstName   *string   `json:"firstname" example:"Farid"`
	LastName    *string   `json:"lastname" example:"Darari"`
	PhoneNumber *string   `json:"phone_number" example:"08123456789"`
	Avatar      *string   `json:"avatar" example:"https://example.com/avatar.png"`
	Point       *int      `json:"point" example:"100"`
}

type ChangePasswordRequest struct {
	OldPassword     string `form:"old_password" json:"old_password" binding:"required"`
	NewPassword     string `form:"new_password" json:"new_password" binding:"required"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" binding:"required"`
}
