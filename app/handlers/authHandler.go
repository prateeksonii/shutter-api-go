package handlers

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/prateeksonii/shutter-api-go/pkg/configs"
	"github.com/prateeksonii/shutter-api-go/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignUp(c *gin.Context) {

	userDto := models.SignUpDto{}

	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.Status(http.StatusBadRequest)
		panic(err)
	}

	existingUser := new(models.User)

	result := configs.Db.Where("username = ?", userDto.Username).First(&existingUser)

	if result.RowsAffected > 0 {
		c.Status(http.StatusConflict)
		panic(errors.New("User already exists"))
	}

	user := new(models.User)
	user.Name = userDto.Name
	user.Username = userDto.Username

	hash, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), 14)

	if err != nil {
		panic(err)
	}

	user.Password = string(hash)

	configs.Db.Create(&user)

	c.JSON(http.StatusCreated, gin.H{
		"ok":   true,
		"user": user,
	})
}

func SignIn(c *gin.Context) {

	userDto := models.SignInDto{}

	if err := c.ShouldBindJSON(userDto); err != nil {
		c.Status(http.StatusBadRequest)
		panic(err)
	}

	user := &models.User{}

	result := configs.Db.Where("username = ?", userDto.Username).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Status(http.StatusNotFound)
		panic(result.Error)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDto.Password)); err != nil {
		c.Status(http.StatusUnauthorized)
		panic(errors.New("invalid username or password"))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.Claims{
		User: *user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(365 * 10 * 24 * time.Hour)),
		},
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))

	if err != nil {
		c.Status(http.StatusInternalServerError)
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":    true,
		"user":  user,
		"token": tokenString,
	})
}

func GetAuthenticatedUser(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	completeUser := &models.User{}

	configs.Db.Model(&models.User{}).Preload("SentInvites").Preload("ReceivedInvites").Where("ID = ?", user.ID).First(&completeUser)

	c.JSON(http.StatusOK, gin.H{
		"ok":   true,
		"user": completeUser,
	})
}

func IsAuthenticated(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if len(authHeader) == 0 {
		c.Status(http.StatusUnauthorized)
		panic(errors.New("no token provided"))
	}

	tokenArray := strings.Split(authHeader, " ")

	if len(tokenArray) < 2 {
		c.Status(http.StatusUnauthorized)
		panic(errors.New("invalid token provided"))
	}

	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenArray[1], claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		c.Status(http.StatusUnauthorized)
		panic(err)
	}

	if !token.Valid {
		c.Status(http.StatusUnauthorized)
		panic(errors.New("invalid token provided"))
	}

	c.Set("user", claims.User)

	c.Next()
}
