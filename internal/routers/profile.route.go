package routers

import (
	"github.com/Darari17/be-tickitz/internal/handlers"
	"github.com/Darari17/be-tickitz/internal/middlewares"
	"github.com/Darari17/be-tickitz/internal/repos"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initProfileRouter(router *gin.Engine, db *pgxpool.Pool) {
	profileRepo := repos.NewProfileRepo(db)
	profileHandler := handlers.NewProfileHandler(profileRepo)

	profile := router.Group("/profile", middlewares.RequiredToken, middlewares.Access("user"))
	profile.GET("", profileHandler.GetProfile)
	profile.PATCH("", profileHandler.UpdateProfile)
	profile.PATCH("/change-password", profileHandler.ChangePassword)
}
