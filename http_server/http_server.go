package httpserver

import (
	"cookbook/http_server/handlers"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

type HttpServer struct {
	app *fiber.App
}

func New(db *gorm.DB) *HttpServer {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Use(swagger.New(
		swagger.Config{
			BasePath: "/",
			FilePath: "./docs/swagger.json",
			Path:     "swagger",
			Title:    "Swagger API Docs",
		},
	))

	app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.SendString("ok")
	})

	handler := handlers.New(db)

	handlers.NewRecepieHandler(app, handler)
	handlers.NewIngredientsHandler(app, handler)
	handlers.NewStepsHandler(app, handler)
	handlers.NewAuthHandlers(app, handler)
	handlers.NewUserHandlers(app, handler)

	return &HttpServer{
		app: app,
	}
}

func (s *HttpServer) Listen(addr string) error {
	return s.app.Listen(addr)
}
