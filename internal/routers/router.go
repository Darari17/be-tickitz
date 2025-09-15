package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	docs "github.com/Darari17/be-tickitz/docs"
	"github.com/Darari17/be-tickitz/internal/middlewares"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(db *pgxpool.Pool, redis *redis.Client) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware)

	initAuthRouter(router, db)
	initMovieRouter(router, db, redis)
	initOrderRouter(router, db)
	initProfileRouter(router, db)
	initAdminRouter(router, db)

	router.Static("/img", "public")

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Rute Salah",
		})
	})

	return router
}
