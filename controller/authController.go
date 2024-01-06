package controller

import (
	"backend/database"
	"backend/model"
	"context"
	"fmt"
	"net/mail"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var SecretKey string

var requireRegister = []string{"name", "email", "password"}
var requireLogin = []string{"email", "password"}

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	SecretKey = os.Getenv("JWTTOKENSECRETKEY")
}

func checkValidInput(require []string, data map[string]string) (string, bool) {
	for i := 0; i < len(require); i++ {
		_, ok := data[require[i]]
		if !ok {
			return "input not complete", false
		}
	}

	return "success", true
}

func isValidMailAddress(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal error",
		})
	}

	message, ok := checkValidInput(requireRegister, data)

	if !ok {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": message,
		})
	}

	if !isValidMailAddress(data["email"]) {
		c.Status(fiber.ErrBadRequest.Code)
		return c.JSON(fiber.Map{
			"message": "Invalid email",
		})
	}

	var foundUser model.User
	filter := bson.D{{Key: "email", Value: data["email"]}}
	err := database.UsersCollection.FindOne(context.TODO(), filter).Decode(&foundUser)
	if err != mongo.ErrNoDocuments {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "email used",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := model.User{
		ID:       primitive.NewObjectID(),
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
		Tasks:    make([]model.Task, 0),
	}

	result, _ := database.UsersCollection.InsertOne(context.TODO(), user)

	return c.JSON(result)
}

func Login(c *fiber.Ctx) error {
	fmt.Println("Getting login request")
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal error",
		})
	}

	message, ok := checkValidInput(requireLogin, data)
	if !ok {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": message,
		})
	}

	if !isValidMailAddress(data["email"]) {
		c.Status(fiber.ErrBadRequest.Code)
		return c.JSON(fiber.Map{
			"message": "Invalid email",
		})
	}

	var user model.User

	filter := bson.D{{Key: "email", Value: data["email"]}}

	err := database.UsersCollection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "user not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal error",
		})
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"]))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.ID.Hex(),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "succress",
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	objectID, err := primitive.ObjectIDFromHex(claims.Issuer)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal error",
		})
	}
	var user model.User
	filter := bson.D{{Key: "_id", Value: objectID}}
	err = database.UsersCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "user not found",
			})
		}
		return err
	}

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func ChangePassword(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	objectID, err := primitive.ObjectIDFromHex(claims.Issuer)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal error",
		})
	}

	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal error",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	update := bson.D{{"$set", bson.D{{"password", password}}}}
	_, err = database.UsersCollection.UpdateByID(context.TODO(), objectID, update)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal error",
		})
	}

	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"message": "success",
	})
}
