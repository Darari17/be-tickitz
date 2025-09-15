package routers

import (
	"github.com/Darari17/be-tickitz/internal/handlers"
	"github.com/Darari17/be-tickitz/internal/repos"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initAuthRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRepo := repos.NewAuthRepo(db)
	authHandler := handlers.NewAuthHandler(authRepo)

	router.POST("/login", authHandler.Login)
	router.POST("/register", authHandler.Register)
}
