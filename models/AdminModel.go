package models

type Admin struct {
	AdminID   string `json:"adminID,omitempty" bson:"adminID,omitempty"`
	SchooldID string `json:"schoolID,omitempty" bson:"schoolID,omitempty"`
	FirstName string `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty" bson:"lastName,omitempty"`
	Email     string `json:"email,omitempty" bson:"email,omitempty"`
	Password  string `json:"password,omitempty" bson:"password,omitempty"`
}