package routers

import (
	"github.com/Darari17/be-tickitz/internal/handlers"
	"github.com/Darari17/be-tickitz/internal/middlewares"
	"github.com/Darari17/be-tickitz/internal/repos"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initAdminRouter(router *gin.Engine, db *pgxpool.Pool) {
	repo := repos.NewAdminRepo(db)
	handler := handlers.NewAdminHandler(repo)

	admin := router.Group("/admin", middlewares.RequiredToken, middlewares.Access("admin"))

	admin.POST("/movies", handler.CreateMovie)
	admin.GET("/movies", handler.GetMovies)
	admin.GET("/movies/:id", handler.GetMovieByID)
	admin.PATCH("/movies/:id", handler.UpdateMovie)
	admin.DELETE("/movies/:id", handler.DeleteMovie)
}
