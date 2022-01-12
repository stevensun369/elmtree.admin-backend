package admin

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
  g := app.Group("/api/admin")

  // admin login
  g.Post("/login", postLogin)

  // get all grades
  g.Get("/grades", AdminMiddleware, getGrades)

  // get subjects
  g.Get("/grades/subjects/:gradeID", AdminMiddleware, getSubjects)

  // here starts the student part of it

  // get students by gradeID
  g.Get("/grades/students/:gradeID", AdminMiddleware, getStudentsByGradeID)

  // get SubjectStudents
  g.Get("/grades/subjects/students/:subjectID", AdminMiddleware, getStudentsBySubjectID)

  // get student
  g.Get("/students/:studentID", AdminMiddleware, getStudent)

  // remove student subject
  g.Put("/students/:studentID/remove/:subjectID", AdminMiddleware, removeStudentSubject)

  // add student subject
  g.Put("/students/:studentID/add", AdminMiddleware, addStudentSubject)


  // here starts the teacher part of it

  // get teacher by cnp
  g.Get("/teachers/:cnp", AdminMiddleware, getTeacher)

  // remove teacher subject
  g.Put("/teachers/:cnp/remove/:subjectID", AdminMiddleware, removeTeacherSubject)

  // add teacher subject
  g.Put("/teachers/:cnp/add", AdminMiddleware, addTeacherSubject)

}