package admin

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
  g := app.Group("/api/admin")

  // admin login
  g.Post("/login", postLogin)

  // get all grades
  g.Get("/grades", adminMiddleware, getGrades)

  // get subjects
  g.Get("/subjects/:gradeID", adminMiddleware, getSubjects)
}