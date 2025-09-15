package routers

import (
	"github.com/Darari17/be-tickitz/internal/handlers"
	"github.com/Darari17/be-tickitz/internal/repos"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func initMovieRouter(router *gin.Engine, db *pgxpool.Pool, redis *redis.Client) {
	movieRepo := repos.NewMovieRepo(db, redis)
	movieHandler := handlers.NewMovieHandler(movieRepo)

	movies := router.Group("/movies")
	movies.GET("/upcoming", movieHandler.GetUpcomingMovies)
	movies.GET("/popular", movieHandler.GetPopularMovies)
	movies.GET("", movieHandler.GetAllMovies)
	movies.GET("/:id", movieHandler.GetMovieDetail)
}
