package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"todo-app/internal/auth"
	"todo-app/internal/errors"
	"todo-app/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	authService auth.Service
}

func NewAuthController(authService auth.Service) *AuthController {
	return &AuthController{authService: authService}
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user  body      model.UserRegister  true  "UserRegister info"
// @Success      201   {object}  model.User
// @Failure      400   {object}  errors.APIError
// @Router       /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var user model.UserRegister
	if err := ctx.ShouldBindJSON(&user); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errorMessages []string
			for _, fieldErr := range validationErrors {
				log.Printf("Validation error: %s", fieldErr.Tag())
				switch fieldErr.Tag() {
				case "required":
					errorMessages = append(errorMessages,
						fmt.Sprintf("%s is required", fieldErr.Field()))
				case "email":
					errorMessages = append(errorMessages,
						fmt.Sprintf("%s must be a valid email", fieldErr.Field()))
				case "min":
					errorMessages = append(errorMessages,
						fmt.Sprintf("%s must be at least %s characters",
							fieldErr.Field(), fieldErr.Param()))
				case "max":
					errorMessages = append(errorMessages,
						fmt.Sprintf("%s must be at most %s characters",
							fieldErr.Field(), fieldErr.Param()))
				default:
					errorMessages = append(errorMessages,
						fmt.Sprintf("%s is invalid", fieldErr.Field()))
				}
			}
			ctx.AbortWithError(400, errors.NewAPIError(
				400,
				"VALIDATION_ERROR",
				strings.Join(errorMessages, ", "),
			))
			return
		}

		ctx.AbortWithError(400, errors.NewAPIErrorWithDetails(
			400,
			"INVALID_PAYLOAD",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	createdUser, err := c.authService.Register(ctx.Request.Context(), &user)
	if err != nil {
		ctx.Error(err)
		return
	}

	createdUser.PasswordHash = ""
	ctx.JSON(201, createdUser)
}

// Login godoc
// @Summary      Authenticate user
// @Description  Logs a user in and returns a JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      model.AuthUser  true  "Login credentials"
// @Success      200          {object}  map[string]string
// @Failure      400          {object}  errors.APIError
// @Failure      401          {object}  errors.APIError
// @Router       /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var authUser model.AuthUser
	if err := ctx.ShouldBindJSON(&authUser); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.NewInvalidIDError())
		return
	}

	token, err := c.authService.Login(ctx.Request.Context(), &authUser)
	if err != nil {
		log.Printf("Login error: %v", err)
		ctx.AbortWithError(http.StatusUnauthorized, errors.NewInvalidCredentialsError())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
