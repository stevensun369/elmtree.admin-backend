package models

import (
	"admin-backend/db"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// the password and the subjectList can be empty
type Student struct {
	StudentID   string           `json:"studentID,omitempty" bson:"studentID,omitempty"`
	FirstName   string           `json:"firstName,omitempty" bson:"firstName,omitempty"`
	DadInitials string           `json:"dadInitials,omitempty" bson:"dadInitials,omitempty"`
	LastName    string           `json:"lastName,omitempty" bson:"lastName,omitempty"`
	CNP         string           `json:"cnp,omitempty" bson:"cnp,omitempty"`
	Password    string           `json:"password" bson:"password"`
	Grade       Grade            `json:"grade" bson:"grade"`
	SubjectList []ShortSubject `json:"subjectList" bson:"subjectList"`
}

// we actually define ShortSubject because it doesn't have the grade field.
type ShortSubject struct {
	SubjectID string `json:"subjectID,omitempty" bson:"subjectID,omitempty"`
	Name      string `json:"name,omitempty" bson:"name,omitempty"`
}

func GetStudentById(studentID string) (Student, error) {
	// getting a student
	var student Student
  studentsColleciton, err := db.GetCollection("students")
  
  studentsColleciton.FindOne(context.Background(), bson.M{"studentID": studentID}).Decode(&student)

	return student, err
}