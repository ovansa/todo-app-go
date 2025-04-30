package integration

import (
	"context"
	"net/http"
	"testing"
	"time"
	"todo-app/internal/auth"
	"todo-app/internal/config"
	"todo-app/internal/controller"
	"todo-app/internal/errors"
	"todo-app/internal/model"
	"todo-app/internal/repository"
	"todo-app/internal/routes"
	"todo-app/pkg/database"
	"todo-app/test"

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
	routes.SetupRoutes(router, suite.authController, nil, suite.authService)
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
		FullName: "Test User",
	}

	w := test.CreateTestRequest(suite.T(), suite.router, "POST", "/auth/register", newUser, "")
	suite.Equal(http.StatusCreated, w.Code)

	var response model.User
	test.ParseResponse(suite.T(), w, &response)

	suite.Equal(newUser.Email, response.Email)
	suite.Empty(response.Password)
	suite.Empty(response.PasswordHash)
}

func (suite *AuthControllerTestSuite) TestRegister_DuplicateEmail() {
	testUser := model.User{
		Email:    "test123@example.com",
		Password: "password123",
		FullName: "Test User",
	}
	err := testUser.HashPassword(suite.authService.GetPepper())
	suite.NoError(err)

	_, err = suite.userRepo.Create(context.Background(), &testUser)
	suite.NoError(err, "Failed to create test user")

	newUser := model.User{
		Email:    "test123@example.com",
		Password: "newpassword123",
		FullName: "New User",
	}
	w := test.CreateTestRequest(suite.T(), suite.router, "POST", "/auth/register", newUser, "")
	suite.Equal(http.StatusConflict, w.Code)

	var response map[string]interface{}
	test.ParseResponse(suite.T(), w, &response)
	suite.Equal(errors.ErrDuplicateEmail.Message, response["message"])
}

func (suite *AuthControllerTestSuite) TestRegister_InvalidInput() {
	tests := []struct {
		name     string
		user     model.User
		expected string
	}{
		{
			name:     "Empty email",
			user:     model.User{Email: "", Password: "password123", FullName: "Test User1"},
			expected: "Email is required",
		},
		{
			name:     "Invalid email",
			user:     model.User{Email: "invalid", Password: "password123", FullName: "Test User2"},
			expected: "Email must be a valid email",
		},
		{
			name:     "Short password",
			user:     model.User{Email: "test@gmail.com", Password: "short", FullName: "Test User3"},
			expected: "Password must be at least 6 characters",
		},
		{
			name:     "Short full name",
			user:     model.User{Email: "test123@gmail.com", Password: "password123", FullName: "a"},
			expected: "FullName must be at least 3 characters",
		},
		{
			name:     "Long full name",
			user:     model.User{Email: "test123@gmail.com", Password: "password123", FullName: "ThisIsALongFullNameInLengthOfMoreThan50CharactersInLength"},
			expected: "FullName must be at most 50 characters",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			w := test.CreateTestRequest(suite.T(), suite.router, "POST", "/auth/register", tt.user, "")
			suite.Equal(http.StatusBadRequest, w.Code)

			var response map[string]interface{}
			test.ParseResponse(suite.T(), w, &response)

			suite.Contains(response["message"], tt.expected)
		})
	}
}

func (suite *AuthControllerTestSuite) TestLogin_Success() {
	user := model.User{
		Email:    "login@example.com",
		Password: "password123",
		FullName: "Test User",
	}

	authUser := model.AuthUser{
		Email:    "login@example.com",
		Password: "password123",
	}

	err := user.HashPassword(suite.authService.GetPepper())
	suite.NoError(err)

	_, err = suite.userRepo.Create(context.Background(), &user)
	suite.NoError(err, "Failed to create test user")

	w := test.CreateTestRequest(suite.T(), suite.router, "POST", "/auth/login", authUser, "")
	suite.Equal(http.StatusOK, w.Code)

	var response map[string]string
	test.ParseResponse(suite.T(), w, &response)
	suite.NotEmpty(response["token"])
}

func (suite *AuthControllerTestSuite) TestLogin_InvalidCredentials() {
	authUser := model.AuthUser{
		Email:    "login@example.com",
		Password: "password123",
	}

	w := test.CreateTestRequest(suite.T(), suite.router, "POST", "/auth/login", authUser, "")
	suite.Equal(http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	test.ParseResponse(suite.T(), w, &response)

	suite.Equal(errors.ErrInvalidCredentials.Message, response["message"])
	suite.Empty(response["token"])
}

func TestAuthControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}
