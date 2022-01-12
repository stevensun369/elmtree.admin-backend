package timetable

import (
	"github.com/gofiber/fiber/v2"

	"admin-backend/admin"
)

func Routes(app *fiber.App) {
  g := app.Group("/api/admin/timetable")

  // getting periods for grade
  g.Get("/:gradeID", admin.AdminMiddleware, getPeriods)

  // modify the period
  g.Put("/:periodID", admin.AdminMiddleware, modifyPeriod)
  
  // unsassign the period
  g.Delete("/:periodID", admin.AdminMiddleware, unassignPeriod)
  
}