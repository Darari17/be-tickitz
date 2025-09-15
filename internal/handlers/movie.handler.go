package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Darari17/be-tickitz/internal/dtos"
	"github.com/Darari17/be-tickitz/internal/repos"
	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	movieRepo *repos.MovieRepo
}

func NewMovieHandler(mr *repos.MovieRepo) *MovieHandler {
	return &MovieHandler{movieRepo: mr}
}

// GetUpcomingMovies godoc
// @Summary Get upcoming movies
// @Description Retrieve paginated list of upcoming movies
// @Tags Movies
// @Produce json
// @Param page query int false "Page number" default(1)
// @Success 200 {object} dtos.SuccessResponse{data=[]models.Movie} "Upcoming movies retrieved successfully"
// @Failure 500 {object} dtos.ErrorResponse "Failed to fetch upcoming movies"
// @Router /movies/upcoming [get]
func (mh *MovieHandler) GetUpcomingMovies(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	movies, err := mh.movieRepo.GetUpcomingMovies(ctx.Request.Context(), page)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch upcoming movies",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    movies,
	})
}

// GetPopularMovies godoc
// @Summary Get popular movies
// @Description Retrieve paginated list of popular movies
// @Tags Movies
// @Produce json
// @Param page query int false "Page number" default(1)
// @Success 200 {object} dtos.SuccessResponse{data=[]models.Movie} "Popular movies retrieved successfully"
// @Failure 500 {object} dtos.ErrorResponse "Failed to fetch popular movies"
// @Router /movies/popular [get]
func (mh *MovieHandler) GetPopularMovies(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	movies, err := mh.movieRepo.GetPopularMovies(ctx.Request.Context(), page)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch popular movies",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    movies,
	})
}

// GetAllMovies godoc
// @Summary Get all movies with filters
// @Description Retrieve paginated list of all movies with optional search and genre filters
// @Tags Movies
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param search query string false "Search query for movie title"
// @Param genre query string false "Genre name filter"
// @Success 200 {object} dtos.SuccessResponse{data=[]models.Movie} "Movies retrieved successfully"
// @Failure 500 {object} dtos.ErrorResponse "Failed to fetch movies"
// @Router /movies [get]
func (mh *MovieHandler) GetAllMovies(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	search := ctx.DefaultQuery("search", "")
	genreName := ctx.DefaultQuery("genre", "")

	movies, err := mh.movieRepo.GetAllMovies(ctx.Request.Context(), page, search, genreName)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch movies",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.SuccessResponse{
		Code:    http.StatusOK,
		Success: true,
		Data:    movies,
	})
}

// GetMovieDetail godoc
// @Summary Get movie detail
// @Description Retrieve detailed information about a specific movie
// @Tags Movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} dtos.SuccessResponse{data=models.Movie} "Movie detail retrieved successfully"
// @Failure 400 {object} dtos.ErrorResponse "Invalid movie ID"
// @Failure 404 {object} dtos.ErrorResponse "Movie not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /movies/{id} [get]
func (mh *MovieHandler) GetMovieDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid movie ID",
		})
		return
	}

	movie, err := mh.movieRepo.GetMovieDetail(ctx.Request.Context(), id)
	if err != nil {
		log.Println(err.Error())
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
