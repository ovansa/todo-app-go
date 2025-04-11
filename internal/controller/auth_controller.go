package controller

import (
	"fmt"
	"net/http"
	"strings"

	"todo-app/internal/auth"
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

func (c *AuthController) Register(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errorMessages []string
			for _, fieldErr := range validationErrors {
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
				default:
					errorMessages = append(errorMessages,
						fmt.Sprintf("%s is invalid", fieldErr.Field()))
				}
			}
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": strings.Join(errorMessages, ", "),
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := c.authService.Register(ctx.Request.Context(), &user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Don't return password hash
	createdUser.PasswordHash = ""
	ctx.JSON(http.StatusCreated, createdUser)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var authUser model.AuthUser
	if err := ctx.ShouldBindJSON(&authUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.authService.Login(ctx.Request.Context(), &authUser)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
