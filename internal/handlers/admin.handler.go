package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Darari17/be-tickitz/internal/dtos"
	"github.com/Darari17/be-tickitz/internal/models"
	"github.com/Darari17/be-tickitz/internal/repos"
	"github.com/Darari17/be-tickitz/internal/utils"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminRepo *repos.AdminRepo
}

func NewAdminHandler(adminRepo *repos.AdminRepo) *AdminHandler {
	return &AdminHandler{adminRepo: adminRepo}
}

// CreateMovie godoc
// @Summary Create a new movie
// @Description Create a new movie with poster & backdrop upload
// @Tags Admin
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Movie title"
// @Param overview formData string false "Movie overview"
// @Param director_name formData string false "Movie director"
// @Param duration formData int true "Movie duration"
// @Param release_date formData string true "Release date (YYYY-MM-DD)"
// @Param popularity formData number false "Movie popularity"
// @Param poster formData file true "Poster image"
// @Param backdrop formData file false "Backdrop image"
// @Param genres formData []string false "Genre names (e.g. Action,Drama or genres=Action&genres=Drama)"
// @Param casts formData []string false "Cast names (e.g. Tom Holland,Zendaya or casts=Tom Holland&casts=Zendaya)"
// @Success 201 {object} dtos.SuccessResponse{data=models.Movie} "Movie created successfully"
// @Failure 400 {object} dtos.ErrorResponse "Invalid request"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /admin/movies [post]
// @Security BearerAuth
func (h *AdminHandler) CreateMovie(ctx *gin.Context) {
	var body dtos.CreateMovieRequest
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	genres := normalizeInputArray(body.Genres)
	casts := normalizeInputArray(body.Casts)

	movie := &models.Movie{
		Title:       body.Title,
		Overview:    body.Overview,
		Director:    body.Director,
		Duration:    body.Duration,
		ReleaseDate: body.ReleaseDate,
		Popularity:  body.Popularity,
	}

	if body.Poster != nil {
		path := utils.SaveImage(ctx, body.Poster, "poster")
		if path == "" {
			return
		}
		movie.Poster = path
	}
	if body.Backdrop != nil {
		path := utils.SaveImage(ctx, body.Backdrop, "backdrop")
		if path == "" {
			return
		}
		movie.Backdrop = path
	}

	var schedules []map[string]interface{}
	for _, s := range body.Schedules {
		date, _ := time.Parse("2006-01-02", s.Date)
		schedules = append(schedules, map[string]interface{}{
			"date":        date,
			"cinema_id":   s.CinemaID,
			"location_id": s.LocationID,
			"time_ids":    s.TimeIDs,
		})
	}

	created, err := h.adminRepo.CreateMovie(ctx, movie, genres, casts, schedules)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, dtos.Response{
		Code:    http.StatusCreated,
		Success: true,
		Message: "Movie created successfully",
		Data:    created,
	})
}

// GetMovies godoc
// @Summary Get all movies (admin)
// @Description Retrieve all movies for admin management
// @Tags Admin
// @Produce json
// @Success 200 {object} dtos.SuccessResponse{data=[]models.Movie} "Movies retrieved successfully"
// @Failure 500 {object} dtos.ErrorResponse "Internal Server Error"
// @Router /admin/movies [get]
// @Security BearerAuth
func (h *AdminHandler) GetMovies(ctx *gin.Context) {
	movies, err := h.adminRepo.GetMovies(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch movies",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    movies,
	})
}

// GetMovieByID godoc
// @Summary Get movie by ID (admin)
// @Description Retrieve a specific movie by ID for admin management
// @Tags Admin
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} dtos.SuccessResponse{data=models.Movie} "Movie retrieved successfully"
// @Failure 404 {object} dtos.ErrorResponse "Movie not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal Server Error"
// @Router /admin/movies/{id} [get]
// @Security BearerAuth
func (h *AdminHandler) GetMovieByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	movie, err := h.adminRepo.GetMovieByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dtos.Response{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "Movie not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    movie,
	})
}

// UpdateMovie godoc
// @Summary Update movie
// @Description Update movie with optional poster & backdrop upload
// @Tags Admin
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Movie ID"
// @Param title formData string false "Movie title"
// @Param overview formData string false "Movie overview"
// @Param director_name formData string false "Movie director"
// @Param duration formData int false "Movie duration"
// @Param release_date formData string false "Release date (YYYY-MM-DD)"
// @Param popularity formData number false "Movie popularity"
// @Param poster formData file false "Poster image"
// @Param backdrop formData file false "Backdrop image"
// @Param genres formData []string false "Genre names (e.g. Action,Drama or genres=Action&genres=Drama)"
// @Param casts formData []string false "Cast names (e.g. Robert Downey Jr,Chris Evans or casts=Robert Downey Jr&casts=Chris Evans)"
// @Success 200 {object} dtos.SuccessResponse{data=models.Movie} "Movie updated successfully"
// @Failure 400 {object} dtos.ErrorResponse "Invalid request"
// @Failure 404 {object} dtos.ErrorResponse "Movie not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /admin/movies/{id} [patch]
// @Security BearerAuth
func (h *AdminHandler) UpdateMovie(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var body dtos.UpdateMovieRequest

	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	_, err := h.adminRepo.GetMovieByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dtos.Response{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "Movie not found",
		})
		return
	}

	genres := normalizeInputArray(body.Genres)
	casts := normalizeInputArray(body.Casts)

	update := make(map[string]interface{})
	if body.Title != nil {
		update["title"] = *body.Title
	}
	if body.Overview != nil {
		update["overview"] = *body.Overview
	}
	if body.Director != nil {
		update["director_name"] = *body.Director
	}
	if body.Duration != nil {
		update["duration"] = *body.Duration
	}
	if body.ReleaseDate != nil {
		update["release_date"] = *body.ReleaseDate
	}
	if body.Popularity != nil {
		update["popularity"] = *body.Popularity
	}
	if body.Poster != nil {
		path := utils.SaveImage(ctx, body.Poster, "poster")
		if path == "" {
			return
		}
		update["poster_path"] = path
	}
	if body.Backdrop != nil {
		path := utils.SaveImage(ctx, body.Backdrop, "backdrop")
		if path == "" {
			return
		}
		update["backdrop_path"] = path
	}

	if err := h.adminRepo.UpdateMovie(ctx, id, update, genres, casts); err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: err.Error(),
		})
		return
	}

	updated, _ := h.adminRepo.GetMovieByID(ctx, id)
	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Movie updated successfully",
		Data:    updated,
	})
}

// DeleteMovie godoc
// @Summary Delete a movie
// @Description Soft delete a movie by ID
// @Tags Admin
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} dtos.SuccessResponse "Movie deleted successfully"
// @Failure 500 {object} dtos.ErrorResponse "Internal Server Error"
// @Router /admin/movies/{id} [delete]
// @Security BearerAuth
func (h *AdminHandler) DeleteMovie(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.adminRepo.SoftDeleteMovie(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to delete movie",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Movie deleted successfully",
	})
}

func normalizeInputArray(input []string) []string {
	var out []string
	for _, item := range input {
		parts := strings.Split(item, ",")
		for _, p := range parts {
			val := strings.TrimSpace(p)
			if val != "" {
				out = append(out, val)
			}
		}
	}
	return out
}
