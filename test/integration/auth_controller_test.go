package integration

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"
	"todo-app/internal/auth"
	"todo-app/internal/config"
	"todo-app/internal/controller"
	"todo-app/internal/model"
	"todo-app/internal/repository"
	"todo-app/pkg/database"
	"todo-app/test/testhelpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

type AuthControllerTestSuite struct {
	suite.Suite
	router         *gin.Engine
	mongoDB        *database.MongoDB
	userRepo       repository.UserRepository
	authService    auth.Service
	authController *controller.AuthController
}

func (suite *AuthControllerTestSuite) SetupSuite() {
	config := config.LoadConfig()

	mongoDB, err := database.NewMongoDB(config.MongoURI, "todo-test-db")
	suite.NoError(err)
	suite.mongoDB = mongoDB

	// Initialize repository and services
	suite.userRepo = repository.NewUserRepository(mongoDB.Database, "users")
	suite.authService = auth.NewAuthService(config.JWTSecret, config.JWTExpiration, config.PasswordPepper, suite.userRepo)
	suite.authController = controller.NewAuthController(suite.authService)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/register", suite.authController.Register)
	router.POST("/auth/login", suite.authController.Login)
	suite.router = router

	// Clear the database before running tests
	_, err = mongoDB.Database.Collection("users").DeleteMany(context.Background(), bson.M{})
	suite.NoError(err)
}

func (suite *AuthControllerTestSuite) TearDownSuite() {
	if suite.mongoDB != nil {
		suite.mongoDB.Close()
	}
}

func (suite *AuthControllerTestSuite) SetupTest() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Drop the entire test database
	err := suite.mongoDB.Database.Drop(ctx)
	suite.Require().NoError(err, "Failed to drop test database")
}

func (suite *AuthControllerTestSuite) TestRegister_Success() {
	newUser := model.User{
		Email:    "test123@example.com",
		Password: "password123",
	}

	w := testhelpers.CreateTestRequest(suite.T(), suite.router, "POST", "/auth/register", newUser, "")
	suite.Equal(http.StatusCreated, w.Code)

	var response model.User
	testhelpers.ParseResponse(suite.T(), w, &response)

	suite.Equal(newUser.Email, response.Email)
	suite.Empty(response.Password)
	suite.Empty(response.PasswordHash)
}

func (suite *AuthControllerTestSuite) TestRegister_DuplicateEmail() {
	testUser := model.User{
		Email:    "test123@example.com",
		Password: "password123",
	}
	err := testUser.HashPassword(suite.authService.GetPepper())
	suite.NoError(err)

	_, err = suite.userRepo.Create(context.Background(), &testUser)
	suite.NoError(err, "Failed to create test user")

	newUser := model.User{
		Email:    testUser.Email,
		Password: "newpassword123",
	}
	w := testhelpers.CreateTestRequest(suite.T(), suite.router, "POST", "/auth/register", newUser, "")

	suite.Equal(http.StatusBadRequest, w.Code)

	var response map[string]string
	testhelpers.ParseResponse(suite.T(), w, &response)
	log.Printf("Response: %v", response)
	suite.Equal("email already exists", response["error"])
}

func (suite *AuthControllerTestSuite) TestRegister_InvalidInput() {
	tests := []struct {
		name     string
		user     model.User
		expected string
	}{
		{
			name:     "Empty email",
			user:     model.User{Email: "", Password: "password123"},
			expected: "Email is required",
		},
		{
			name:     "Invalid email",
			user:     model.User{Email: "invalid", Password: "password123"},
			expected: "Email must be a valid email",
		},
		{
			name:     "Short password",
			user:     model.User{Email: "test@gmail.com", Password: "short"},
			expected: "Password must be at least 6 characters",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			w := testhelpers.CreateTestRequest(suite.T(), suite.router, "POST", "/auth/register", tt.user, "")
			suite.Equal(http.StatusBadRequest, w.Code)

			var response map[string]string
			testhelpers.ParseResponse(suite.T(), w, &response)

			suite.Contains(response["error"], tt.expected)
		})
	}
}

func TestAuthControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}
