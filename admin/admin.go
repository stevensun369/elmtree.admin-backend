package admin

import (
	"github.com/gofiber/fiber/v2"

	// internal packages
	"admin-backend/db"
	"admin-backend/models"

	// std
	"context"
	"encoding/json"
	"fmt"

	// mongo
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getGrades(c *fiber.Ctx) error {
  schoolIDLocals := fmt.Sprintf("%v", c.Locals("schoolID"))
  var schoolID string
  json.Unmarshal([]byte(schoolIDLocals), &schoolID)

  var grades []models.Grade
  gradesCollection, err := db.GetCollection("grades")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  options := options.Find()
  options.SetSort(bson.D{{Key: "gradeNumber", Value: 1}, {Key: "gradeLetter", Value: 1}})
  cursor, err := gradesCollection.Find(context.Background(), bson.M{
    "schoolID": schoolID,
  }, options)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if err = cursor.All(context.Background(), &grades); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if len(grades) == 0 {
    grades = []models.Grade {}
  }

  return c.JSON(grades)
}

func getSubjects(c *fiber.Ctx) error {

  gradeID := c.Params("gradeID")

  schoolIDLocals := fmt.Sprintf("%v", c.Locals("schoolID"))
  var schoolID string
  json.Unmarshal([]byte(schoolIDLocals), &schoolID)

  var subjects []models.Subject
  subjectsCollection, err := db.GetCollection("subjects")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  cursor, err := subjectsCollection.Find(context.Background(), bson.M{
    "grade.gradeID": gradeID,
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if err = cursor.All(context.Background(), &subjects); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if len(subjects) == 0 {
    subjects = []models.Subject {}
  }

  return c.JSON(subjects)
}