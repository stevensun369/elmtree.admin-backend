package admin

import (
	"github.com/gofiber/fiber/v2"

	// internal packages
	"admin-backend/db"
	"admin-backend/models"
	"admin-backend/utils"

	// std
	"context"
	"encoding/json"
	"strconv"

	// mongo
	"go.mongodb.org/mongo-driver/bson"
)

// @desc   Get grades for students by grade id
// @route  GET /api/admin/grades/students/:gradeID
// @access Private
func getStudentsByGradeID(c *fiber.Ctx) error {
	gradeID := c.Params("gradeID")

	students, err := db.GetStudents(bson.M{
		"grade.gradeID": gradeID,
	}, db.EmptySort)
	utils.CheckError(c, err)

	if len(students) == 0 {
		students = []models.Student{}
	}

	return c.JSON(students)
}


// @desc   Get grades for students by subjectID
// @route  GET /api/admin/grades/subjects/students/:subjectID
// @access Private
func getStudentsBySubjectID(c *fiber.Ctx) error {
	subjectID := c.Params("subjectID")

	students, err := db.GetStudents(bson.M{
		"subject.subjectID": subjectID,
	}, db.EmptySort)
	utils.CheckError(c, err)

	if len(students) == 0 {
		students = []models.Student{}
	}

	return c.JSON(students)
}

// @desc   Get grades for students by subjectID
// @route  GET /api/admin/students/:studentID
// @access Private
func getStudent(c *fiber.Ctx) error {
	studentID := c.Params("studentID")

	student, err := db.GetStudentByID(studentID)
	utils.CheckMessageError(c, err, "Nu există niciun elev cu ID-ul de elev introdus.")

	return c.JSON(student)
}
// @desc   remove subject of student
// @route  PUT /api/admin/students/:studentID/remove/:subjectID
// @access Private
func removeStudentSubject(c *fiber.Ctx) error {
	studentID := c.Params("studentID")
	subjectID := c.Params("subjectID")

	student, err := db.GetStudentByID(studentID)
	utils.CheckMessageError(c, err, "Nu există niciun elev cu ID-ul de elev introdus.")

	// get the student subject list, and eliminate the subject
	oldSubjectList := student.SubjectList
	var newSubjectList []models.ShortSubject
	for _, subject := range oldSubjectList {
		if subject.SubjectID != subjectID {
			newSubjectList = append(newSubjectList, subject)
		}
	}

	db.Students.FindOneAndUpdate(context.Background(), bson.M{
		"studentID": studentID,
	}, bson.M{
		"$set": bson.M{"subjectList": newSubjectList},
	})

	return c.JSON(newSubjectList)
}

// @desc   add subject of student
// @route  PUT /api/admin/students/:studentID/add
// @access Private
func addStudentSubject(c *fiber.Ctx) error {
	// getting body and unmarshalling it into a body variable
	var body map[string]string
	json.Unmarshal(c.Body(), &body)

	subjectName := body["name"]
	subjectGradeNumber, _ := strconv.Atoi(body["gradeNumber"])
	subjectGradeLetter := body["gradeLetter"]

	studentID := c.Params("studentID")

	student, err := db.GetStudentByID(studentID)
	utils.CheckMessageError(c, err, "Nu există niciun elev cu ID-ul de elev introdus.")

	// subject
	var subject models.Subject
	if err := db.Subjects.FindOne(context.Background(), bson.M{
		"name":              subjectName,
		"grade.gradeNumber": subjectGradeNumber,
		"grade.gradeLetter": subjectGradeLetter,
	}).Decode(&subject); err != nil {
		return utils.MessageError(c, "Nu există niciun subiect cu datele acestea.")
	}

	// new subjectList
	shortSubject := models.ShortSubject{
		Name:      subject.Name,
		SubjectID: subject.SubjectID,
	}
	oldSubjectList := student.SubjectList
	newSubjectList := append(oldSubjectList, shortSubject)

	db.Students.FindOneAndUpdate(context.Background(), bson.M{
		"studentID": studentID,
	}, bson.M{
		"$set": bson.M{"subjectList": newSubjectList},
	})

	return c.JSON(newSubjectList)
}