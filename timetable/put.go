package timetable

import (
	"admin-backend/db"
	"admin-backend/models"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func modifyPeriod(c *fiber.Ctx) error {
	periodID := c.Params("periodID")
	periodsCollection, err := db.GetCollection("periods")
	if err != nil {
	  return c.Status(500).SendString(fmt.Sprintf("%v", err))
	}

	// the body
	var body map[string]string
	json.Unmarshal(c.Body(), &body)

	room := body["room"]
	subjectName := body["subjectName"]
	subjectGradeNumber, _ := strconv.Atoi(body["subjectGradeNumber"])
	subjectGradeLetter := body["subjectGradeLetter"]

	// getting the subject
	subjectsCollection, err := db.GetCollection("subjects")
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("%v", err))
	}

	var subject models.Subject
	err = subjectsCollection.FindOne(context.Background(), bson.M{
		"name":              subjectName,
		"grade.gradeNumber": subjectGradeNumber,
		"grade.gradeLetter": subjectGradeLetter,
	}).Decode(&subject)
	if err != nil {
		return c.Status(500).JSON(bson.M{
			"message": "Nu exista nicio materie cu numele respectiv la aceasta clasa.",
		})
	}

  shortSubject := models.ShortSubject {
    Name: subject.Name,
    SubjectID: subject.SubjectID,
  }

	// getting the teacher
	teachersCollection, err := db.GetCollection("teachers")
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("%v", err))
	}

	var teacher models.Teacher
	err = teachersCollection.FindOne(context.Background(), bson.M{
		"subjectList.name": subjectName,
		"subjectList.grade.gradeNumber": subjectGradeNumber,
		"subjectList.grade.gradeLetter": subjectGradeLetter,
	}).Decode(&teacher)
	if err != nil {
		return c.Status(500).JSON(bson.M{
			"message": "Nu exista niciun profesor cu datele respective.",
		})
	}

	shortTeacher := models.ShortTeacher {
		TeacherID: teacher.TeacherID,
		FirstName: teacher.FirstName,
		LastName: teacher.LastName,
	}

  var period models.Period
  err = periodsCollection.FindOneAndUpdate(context.Background(), bson.M{
    "periodID": periodID,
  }, bson.M{
    "$set": bson.M{
      "subject": shortSubject,
			"teacher": shortTeacher,
      "room": room,
      "assigned": true,
    },
  }).Decode(&period)
  if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("%v", err))
	}
  period.Subject = shortSubject
	period.Teacher = shortTeacher
  period.Room = room
  period.Assigned = true

	return c.JSON(period)
}

func unassignPeriod(c *fiber.Ctx) error {
  periodID := c.Params("periodID")
  periodsCollection, err := db.GetCollection("periods")
	if err != nil {
	  return c.Status(500).SendString(fmt.Sprintf("%v", err))
	}
  var period models.Period
  err = periodsCollection.FindOneAndUpdate(context.Background(), bson.M{
    "periodID": periodID,
  }, bson.M{
    "$set": bson.M{
      "subject": models.ShortSubject{},
			"teacher": models.ShortTeacher{},
      "room": "",
      "assigned": false,
    },
  }).Decode(&period)
  if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("%v", err))
	}
  period.Subject = models.ShortSubject{}
	period.Teacher = models.ShortTeacher {}
  period.Room = ""
  period.Assigned = false

  return c.JSON(period)
}