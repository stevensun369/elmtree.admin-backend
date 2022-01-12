package admin

import (
	"admin-backend/db"
	"admin-backend/models"
	"admin-backend/utils"
	"context"
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// @desc   get teacher
// @route  GET /api/admin/teacher/:cnp
// @access Private
func getTeacher(c *fiber.Ctx) error {
	cnp := c.Params("cnp")

	teacher, err := db.GetTeacherByCNP(cnp)
	utils.CheckMessageError(c, err, "Nu există niciun profesor cu CNP-ul introdus.")

	return c.JSON(teacher)
}

// @desc   remove subject of teacher
// @route  PUT /api/admin/teacher/:cnp/remove/:subjectID
// @access Private
func removeTeacherSubject(c *fiber.Ctx) error {
	cnp := c.Params("cnp")
	subjectID := c.Params("subjectID")

	teacher, err := db.GetTeacherByCNP(cnp)
	utils.CheckMessageError(c, err, "Nu există niciun profesor cu CNP-ul introdus.")

	// get the teacher subject list, and eliminate the subject
	oldSubjectList := teacher.SubjectList
	var newSubjectList []models.Subject
	for _, subject := range oldSubjectList {
		if subject.SubjectID != subjectID {
			newSubjectList = append(newSubjectList, subject)
		}
	}

	db.Teachers.FindOneAndUpdate(context.Background(), bson.M{
		"cnp": cnp,
	}, bson.M{
		"$set": bson.M{"subjectList": newSubjectList},
	})

	return c.JSON(newSubjectList)
}

// @desc   add subject of teacher
// @route  PUT /api/admin/teacher/:cnp/add/
// @access Private
func addTeacherSubject(c *fiber.Ctx) error {
	// getting body and unmarshalling it into a body variable
	var body map[string]string
	json.Unmarshal(c.Body(), &body)

	subjectName := body["name"]
	subjectGradeNumber, _ := strconv.Atoi(body["gradeNumber"])
	subjectGradeLetter := body["gradeLetter"]

	cnp := c.Params("cnp")

	teacher, err := db.GetTeacherByCNP(cnp)
	utils.CheckMessageError(c, err, "Nu există niciun profesor cu CNP-ul introdus.")

	
	var subject models.Subject
	if err = db.Subjects.FindOne(context.Background(), bson.M{
		"name":              subjectName,
		"grade.gradeNumber": subjectGradeNumber,
		"grade.gradeLetter": subjectGradeLetter,
	}).Decode(&subject); err != nil {
		return utils.MessageError(c, "Nu există niciun subiect cu datele acestea.")
	}

	oldSubjectList := teacher.SubjectList
	newSubjectList := append(oldSubjectList, subject)

	db.Teachers.FindOneAndUpdate(context.Background(), bson.M{
		"cnp": cnp,
	}, bson.M{
		"$set": bson.M{"subjectList": newSubjectList},
	})

	return c.JSON(newSubjectList)
}