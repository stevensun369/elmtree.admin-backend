package timetable

import (
	"admin-backend/db"
	"admin-backend/models"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getPeriods(c *fiber.Ctx) error {
  gradeID := c.Params("gradeID")
  periodsCollection, err := db.GetCollection("periods")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  options := options.Find()
  options.SetSort(bson.D{{Key: "day", Value: 1}, {Key: "interval", Value: 1}})
  var periods []models.Period
  cursor, err := periodsCollection.Find(context.Background(), bson.M{
    "grade.gradeID": gradeID,
  }, options)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &periods); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(periods)
}

