package controller

import (
	"backend/database"
	"backend/model"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddTask(c *fiber.Ctx) error {
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
		return err
	}

	var newTask model.Task

	if err := c.BodyParser(&newTask); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal error",
		})
	}

	id := primitive.NewObjectID()

	newTask.ID = id
	if newTask.Priority < 0 {
		newTask.Priority = 0
	} else if newTask.Priority > 3 {
		newTask.Priority = 3
	}

	if newTask.Tags == nil {
		newTask.Tags = make([]string, 0)
	}

	// filter := bson.D{{Key: "_id", Value: objectID}}
	update := bson.D{{"$push", bson.D{{"tasks", newTask}}}}

	_, err = database.UsersCollection.UpdateByID(context.TODO(), objectID, update)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal error",
		})
	}

	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"message": "GOOD",
		"_id":     id.Hex(),
	})
}

func DeleteTask(c *fiber.Ctx) error {
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

	userObjectID, err := primitive.ObjectIDFromHex(claims.Issuer)
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

	itemObjectId, err := primitive.ObjectIDFromHex(data["_id"])
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal error",
		})
	}

	update := bson.D{{"$pull", bson.D{{"tasks", bson.D{{"_id", itemObjectId}}}}}}

	res, err := database.UsersCollection.UpdateByID(context.TODO(), userObjectID, update)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "pull error",
		})
	}

	c.Status(fiber.StatusOK)
	return c.JSON(res)
}
