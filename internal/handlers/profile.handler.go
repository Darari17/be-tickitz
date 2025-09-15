package handlers

import (
	"net/http"

	"github.com/Darari17/be-tickitz/internal/dtos"
	"github.com/Darari17/be-tickitz/internal/models"
	"github.com/Darari17/be-tickitz/internal/repos"
	"github.com/Darari17/be-tickitz/internal/utils"
	"github.com/Darari17/be-tickitz/pkg"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileRepo *repos.ProfileRepo
}

func NewProfileHandler(pr *repos.ProfileRepo) *ProfileHandler {
	return &ProfileHandler{profileRepo: pr}
}

// GetProfile godoc
// @Summary Get user profile
// @Description Retrieve profile information for the authenticated user
// @Tags Profile
// @Produce json
// @Success 200 {object} dtos.SuccessResponse{data=dtos.ProfileResponse} "Profile retrieved successfully"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure 404 {object} dtos.ErrorResponse "Profile not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /profile [get]
// @Security BearerAuth
func (ph *ProfileHandler) GetProfile(ctx *gin.Context) {
	userID, _, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized: " + err.Error(),
		})
		return
	}

	profile, err := ph.profileRepo.GetProfile(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dtos.Response{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "Profile not found",
		})
		return
	}

	res := dtos.ProfileResponse{
		UserID:      profile.UserID,
		FirstName:   profile.FirstName,
		LastName:    profile.LastName,
		PhoneNumber: profile.PhoneNumber,
		Avatar:      profile.Avatar,
		Point:       profile.Point,
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    res,
	})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update profile information for the authenticated user
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Param firstname formData string false "First name"
// @Param lastname formData string false "Last name"
// @Param phone_number formData string false "Phone number"
// @Param avatar formData file false "Avatar image"
// @Success 200 {object} dtos.Response "Profile updated successfully"
// @Failure 400 {object} dtos.Response "Invalid request body"
// @Failure 401 {object} dtos.Response "Unauthorized"
// @Failure 500 {object} dtos.Response "Failed to update profile"
// @Router /profile [patch]
// @Security BearerAuth
func (ph *ProfileHandler) UpdateProfile(ctx *gin.Context) {
	userID, _, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized: " + err.Error(),
		})
		return
	}

	var req dtos.UpdateProfileRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	var avatarPath *string
	if req.Avatar != nil {
		filename := utils.SaveImage(ctx, req.Avatar, "avatars")
		avatarPath = &filename
	}

	profile := &models.Profile{
		UserID:      userID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Avatar:      avatarPath,
	}

	if err := ph.profileRepo.UpdateProfile(ctx.Request.Context(), profile); err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to update profile",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Profile updated successfully",
	})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change password using old, new, and confirm password
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Param old_password formData string true "Old password"
// @Param new_password formData string true "New password"
// @Param confirm_password formData string true "Confirm password"
// @Success 200 {object} dtos.Response "Password changed successfully"
// @Failure 400 {object} dtos.Response "Invalid request body or password mismatch"
// @Failure 401 {object} dtos.Response "Unauthorized or invalid old password"
// @Failure 500 {object} dtos.Response "Failed to update password"
// @Router /profile/change-password [patch]
// @Security BearerAuth
func (ph *ProfileHandler) ChangePassword(ctx *gin.Context) {
	userID, _, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized: " + err.Error(),
		})
		return
	}

	var req dtos.ChangePasswordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "New password and confirm password do not match",
		})
		return
	}

	hashedPassword, err := ph.profileRepo.VerifyPassword(ctx.Request.Context(), userID, req.OldPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Invalid old password",
		})
		return
	}

	hashConfig := pkg.NewHashConfig()
	ok, err := hashConfig.CompareHashAndPassword(req.OldPassword, hashedPassword)
	if err != nil || !ok {
		ctx.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Invalid old password",
		})
		return
	}

	hashConfig.UseRecommended()
	newHashed, err := hashConfig.GenHash(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to hash new password",
		})
		return
	}

	if err := ph.profileRepo.UpdatePassword(ctx.Request.Context(), userID, newHashed); err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to update password",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Password changed successfully",
	})
}
