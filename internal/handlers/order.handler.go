package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Darari17/be-tickitz/internal/dtos"
	"github.com/Darari17/be-tickitz/internal/models"
	"github.com/Darari17/be-tickitz/internal/repos"
	"github.com/Darari17/be-tickitz/internal/utils"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderRepo *repos.OrderRepo
}

func NewOrderHandler(or *repos.OrderRepo) *OrderHandler {
	return &OrderHandler{orderRepo: or}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new movie ticket order
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body dtos.CreateOrderRequest true "Order creation data"
// @Success 201 {object} dtos.Response{data=models.Order} "Order created successfully"
// @Failure 400 {object} dtos.Response "Invalid request payload or seat codes"
// @Failure 401 {object} dtos.Response "Unauthorized"
// @Failure 500 {object} dtos.Response "Failed to create order"
// @Router /orders [post]
// @Security BearerAuth
func (oh *OrderHandler) CreateOrder(ctx *gin.Context) {
	var req dtos.CreateOrderRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	userID, _, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	seatIDs, err := oh.orderRepo.GetSeatIDsByCodes(ctx.Request.Context(), req.SeatCodes)
	if err != nil || len(seatIDs) != len(req.SeatCodes) {
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid seat codes",
		})
		return
	}

	order := &models.Order{
		UserID:     userID,
		ScheduleID: req.ScheduleID,
		PaymentID:  req.PaymentID,
		FullName:   req.FullName,
		Email:      req.Email,
		Phone:      req.Phone,
	}

	newOrder, err := oh.orderRepo.CreateOrder(ctx.Request.Context(), order, seatIDs)
	if err != nil {
		log.Println("CreateOrder error:", err)
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to create order",
		})
		return
	}

	ctx.JSON(http.StatusCreated, dtos.Response{
		Code:    http.StatusCreated,
		Success: true,
		Data:    newOrder,
	})
}

// GetSchedules godoc
// @Summary Get movie schedules
// @Description Retrieve all schedules for a specific movie
// @Tags Orders
// @Produce json
// @Param movie_id query int true "Movie ID"
// @Success 200 {object} dtos.Response{data=[]models.Schedule} "Schedules retrieved successfully"
// @Failure 400 {object} dtos.Response "Invalid movie_id"
// @Failure 500 {object} dtos.Response "Failed to fetch schedules"
// @Router /orders/schedules [get]
// @Security BearerAuth
func (oh *OrderHandler) GetSchedules(ctx *gin.Context) {
	movieIDStr := ctx.Query("movie_id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid movie_id",
		})
		return
	}

	schedules, err := oh.orderRepo.GetSchedules(ctx.Request.Context(), movieID)
	if err != nil {
		log.Println("GetSchedules error:", err)
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch schedules",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    schedules,
	})
}

// GetAvailableSeats godoc
// @Summary Get available seats
// @Description Retrieve available seats for a specific schedule
// @Tags Orders
// @Produce json
// @Param schedule_id query int true "Schedule ID"
// @Success 200 {object} dtos.Response{data=[]models.Seat} "Available seats retrieved successfully"
// @Failure 400 {object} dtos.Response "Invalid schedule_id"
// @Failure 500 {object} dtos.Response "Failed to fetch seats"
// @Router /orders/seats [get]
// @Security BearerAuth
func (oh *OrderHandler) GetAvailableSeats(ctx *gin.Context) {
	scheduleIDStr := ctx.Query("schedule_id")
	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid schedule_id",
		})
		return
	}

	seats, err := oh.orderRepo.GetAvailableSeats(ctx.Request.Context(), scheduleID)
	if err != nil {
		log.Println("GetAvailableSeats error:", err)
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch seats",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    seats,
	})
}

// GetTransactionDetail godoc
// @Summary Get transaction detail
// @Description Retrieve detailed information about a specific transaction
// @Tags Orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} dtos.Response{data=models.OrderDetail} "Transaction detail retrieved successfully"
// @Failure 400 {object} dtos.Response "Invalid order ID"
// @Failure 404 {object} dtos.Response "Order not found"
// @Failure 500 {object} dtos.Response "Internal server error"
// @Router /orders/{id} [get]
// @Security BearerAuth
func (oh *OrderHandler) GetTransactionDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid order ID",
		})
		return
	}

	order, err := oh.orderRepo.GetTransactionDetail(ctx.Request.Context(), id)
	if err != nil {
		log.Println("GetTransactionDetail error:", err)
		ctx.JSON(http.StatusNotFound, dtos.Response{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "Order not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    order,
	})
}

// GetOrderHistory godoc
// @Summary Get user order history
// @Description Retrieve order history for the authenticated user
// @Tags Orders
// @Produce json
// @Success 200 {object} dtos.Response{data=[]models.OrderDetail} "Order history retrieved successfully"
// @Failure 401 {object} dtos.Response "Unauthorized"
// @Failure 500 {object} dtos.Response "Failed to fetch order history"
// @Router /orders/history [get]
// @Security BearerAuth
func (oh *OrderHandler) GetOrderHistory(ctx *gin.Context) {
	userID, _, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	orders, err := oh.orderRepo.GetOrderHistory(ctx.Request.Context(), userID)
	if err != nil {
		log.Println("GetOrderHistory error:", err)
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch order history",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    orders,
	})
}
