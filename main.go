package main

import (

	// gofiber
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// internal pacakges
	"admin-backend/admin"
	"admin-backend/db"
)

func main() {
  app := fiber.New()

  app.Use(cors.New(cors.Config{
    AllowOrigins: "*",
  }))

  db.InitDatabase()
  admin.Routes(app)

  app.Get("/", func (c *fiber.Ctx) error {
    return c.SendString("api is running")
  })

  log.Fatal(app.Listen(":9999"))
}