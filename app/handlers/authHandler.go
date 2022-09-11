package handlers

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/prateeksonii/shutter-api-go/pkg/configs"
	"github.com/prateeksonii/shutter-api-go/pkg/models"
	"github.com/prateeksonii/shutter-api-go/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignUp(c *fiber.Ctx) error {

	userDto := &models.SignUpDto{}

	if err := c.BodyParser(userDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := utils.NewValidator()

	if err := validate.Struct(userDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	existingUser := new(models.User)

	result := configs.Db.Where("username = ?", userDto.Username).First(&existingUser)

	if result.RowsAffected > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "User already exists",
		})
	}

	user := new(models.User)
	user.Name = userDto.Name
	user.Username = userDto.Username

	hash, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), 14)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	user.Password = string(hash)

	configs.Db.Create(&user)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"ok":   true,
		"user": user,
	})
}

func SignIn(c *fiber.Ctx) error {

	userDto := &models.SignInDto{}

	if err := c.BodyParser(userDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	validate := utils.NewValidator()

	if err := validate.Struct(userDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	user := &models.User{}

	result := configs.Db.Where("username = ?", userDto.Username).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": result.Error.Error(),
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDto.Password)); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return errors.New("invalid username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.Claims{
		User: *user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(365 * 10 * 24 * time.Hour)),
		},
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"ok":    true,
		"user":  user,
		"token": tokenString,
	})
}

func GetAuthenticatedUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)

	if !ok {
		c.Status(fiber.StatusUnauthorized)
		return errors.New("invalid user")
	}

	return c.JSON(fiber.Map{
		"ok":   true,
		"user": user,
	})
}

func IsAuthenticated(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if len(authHeader) == 0 {
		c.Status(fiber.StatusUnauthorized)
		return errors.New("no token provided")
	}

	tokenArray := strings.Split(authHeader, " ")

	if len(tokenArray) < 2 {
		c.Status(fiber.StatusUnauthorized)
		return errors.New("invalid token provided")
	}

	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenArray[1], claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return err
	}

	if !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return errors.New("invalid token provided")
	}

	c.Locals("user", claims.User)

	return c.Next()
}
