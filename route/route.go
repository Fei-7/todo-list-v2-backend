package route

import (
	"backend/controller"

	"github.com/gofiber/fiber/v2"
)

func SetUp(app *fiber.App) {

	app.Post("/api/register", controller.Register)
	app.Post("/api/login", controller.Login)
	app.Get("/api/user", controller.User)
	app.Post("/api/logout", controller.Logout)
	app.Post("/api/password", controller.ChangePassword)

	app.Post("/api/task", controller.AddTask)
	app.Delete("/api/task", controller.DeleteTask)

}
