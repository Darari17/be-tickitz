package routers

import (
	"github.com/Darari17/be-tickitz/internal/handlers"
	"github.com/Darari17/be-tickitz/internal/middlewares"
	"github.com/Darari17/be-tickitz/internal/repos"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initOrderRouter(router *gin.Engine, db *pgxpool.Pool) {
	orderRepo := repos.NewOrderRepo(db)
	orderHandler := handlers.NewOrderHandler(orderRepo)

	orderGroup := router.Group("/orders", middlewares.RequiredToken, middlewares.Access("admin", "user"))
	orderGroup.POST("", orderHandler.CreateOrder)
	orderGroup.GET("/history", orderHandler.GetOrderHistory)
	orderGroup.GET("/schedules", orderHandler.GetSchedules)
	orderGroup.GET("/seats", orderHandler.GetAvailableSeats)
	orderGroup.GET("/:id", orderHandler.GetTransactionDetail)
}
