package handlers

import (
	"net/http"
	"time"

	"github.com/ZAUakaAlexey/backend_go/internal/config"
	"github.com/ZAUakaAlexey/backend_go/internal/database"
	"github.com/ZAUakaAlexey/backend_go/internal/models"
	"github.com/ZAUakaAlexey/backend_go/internal/validators"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type SignupInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,strongpassword"`
	Name     string `json:"name" binding:"required,fullname"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileInput struct {
	Name     string `json:"name" binding:"omitempty,fullname"`
	Phone    string `json:"phone" binding:"omitempty,phone"`
	Username string `json:"username" binding:"omitempty,username"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,strongpassword"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func Signup(c *gin.Context) {
	var input SignupInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": validators.FormatValidationErrors(err),
		})
		return
	}

	var existingUser models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this email already exists",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := models.User{
		Email:    input.Email,
		Password: string(hashedPassword),
		Name:     input.Name,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User:  user,
	})
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": validators.FormatValidationErrors(err),
		})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}

func generateToken(userID uint) (string, error) {
	cfg, _ := config.LoadConfig()

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}
