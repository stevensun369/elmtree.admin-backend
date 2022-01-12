package admin

import (
	"github.com/gofiber/fiber/v2"

	// internal packages
	"admin-backend/db"
	"admin-backend/models"
	"admin-backend/utils"

	// std

	"encoding/json"
	"fmt"

	// mongo
	"go.mongodb.org/mongo-driver/bson"
)

// @desc   Get grades
// @route  GET /api/admin/grades
// @access Private
func getGrades(c *fiber.Ctx) error {
  schoolIDLocals := fmt.Sprintf("%v", c.Locals("schoolID"))
  var schoolID string
  json.Unmarshal([]byte(schoolIDLocals), &schoolID)

  grades, err := db.GetGrades(bson.M{
    "schoolID": schoolID,
  }, db.GradeSort)
  utils.CheckError(c, err) 

  if len(grades) == 0 {
    grades = []models.Grade {}
  }

  return c.JSON(grades)
}

// @desc   Get subjects
// @route  GET /api/admin/subjects
// @access Private
func getSubjects(c *fiber.Ctx) error {

  gradeID := c.Params("gradeID")

  subjects, err := db.GetSubjects(bson.M{
    "grade.gradeID": gradeID,
  }, db.GradeSort)
  utils.CheckError(c, err) 

  if len(subjects) == 0 {
    subjects = []models.Subject {}
  }

  return c.JSON(subjects)
}
