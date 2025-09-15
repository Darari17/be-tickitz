package handlers

import (
	"log"
	"net/http"

	"github.com/Darari17/be-tickitz/internal/dtos"
	"github.com/Darari17/be-tickitz/internal/models"
	"github.com/Darari17/be-tickitz/internal/repos"
	"github.com/Darari17/be-tickitz/pkg"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authRepo *repos.AuthRepo
}

func NewAuthHandler(ar *repos.AuthRepo) *AuthHandler {
	return &AuthHandler{
		authRepo: ar,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body dtos.AuthRequest true "Login credentials"
// @Success 200 {object} dtos.SuccessResponse{data=dtos.AuthResponse} "Login successful"
// @Failure 400 {object} dtos.ErrorResponse "Invalid request body"
// @Failure 401 {object} dtos.ErrorResponse "Invalid email or password"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /login [post]
func (ah *AuthHandler) Login(ctx *gin.Context) {
	body := dtos.AuthRequest{}
	if err := ctx.ShouldBind(&body); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid Request Body",
		})
		return
	}

	user, err := ah.authRepo.FindByEmail(ctx.Request.Context(), body.Email)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Something went wrong",
		})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Invalid Email or Password",
		})
		return
	}

	var hash pkg.HashConfig
	valid, err := hash.CompareHashAndPassword(body.Password, user.Password)
	if err != nil || !valid {
		log.Println(err.Error())
		ctx.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Invalid Email or Password",
		})
		return
	}

	claim := pkg.NewJWTClaims(user.ID, string(user.Role))
	token, err := claim.GenToken()
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to Generate Token",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data: dtos.AuthResponse{
			UserID: user.ID,
			Token:  token,
		},
	})
}

// Register godoc
// @Summary User registration
// @Description Register a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body dtos.AuthRequest true "Registration data"
// @Success 201 {object} dtos.SuccessResponse "Registration successful"
// @Failure 400 {object} dtos.ErrorResponse "Invalid request body"
// @Failure 409 {object} dtos.ErrorResponse "Email already exists"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /register [post]
func (ah *AuthHandler) Register(ctx *gin.Context) {
	body := dtos.AuthRequest{}
	if err := ctx.ShouldBind(&body); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid Request Body",
		})
		return
	}

	var hash pkg.HashConfig
	hash.UseRecommended()
	hashed, err := hash.GenHash(body.Password)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to Hash Password",
		})
		return
	}

	user := models.User{
		Email:    body.Email,
		Password: hashed,
	}

	err = ah.authRepo.CreateUser(ctx.Request.Context(), &user)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusConflict, dtos.Response{
			Code:    http.StatusConflict,
			Success: false,
			Message: "Email Already Exists",
		})
		return
	}

	ctx.JSON(http.StatusCreated, dtos.Response{
		Code:    http.StatusCreated,
		Success: true,
		Message: "Register Successfully",
	})
}
