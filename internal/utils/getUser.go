package utils

import (
	"errors"

	"github.com/Darari17/be-tickitz/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserFromContext(ctx *gin.Context) (uuid.UUID, string, error) {
	claims, exists := ctx.Get("claims")
	if !exists {
		return uuid.Nil, "", errors.New("claims not found in context - token might be missing")
	}

	userClaims, ok := claims.(*pkg.Claims)
	if !ok {
		return uuid.Nil, "", errors.New("invalid claims format")
	}

	return userClaims.UserID, userClaims.Role, nil
}
