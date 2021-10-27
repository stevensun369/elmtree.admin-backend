package admin

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"

	// internal packages
	"admin-backend/db"
	"admin-backend/models"
	"admin-backend/utils"

	// std
	"context"
	"encoding/json"
	"fmt"

	// mongo
	"go.mongodb.org/mongo-driver/bson"

	// bcrypt
	"golang.org/x/crypto/bcrypt"
)

func adminMiddleware(c *fiber.Ctx) error {
  var token string

  authHeader := c.Get("Authorization")

  if string(authHeader) != "" && strings.HasPrefix(string(authHeader), "Bearer") {
    token = strings.Fields(string(authHeader))[1]

    // we're just parsing the token: maybe I will put it in the utils
    claims := &utils.AdminClaims{}
    tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface {}, error) {
      return utils.JWTKey, nil
    })

    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    if !tkn.Valid {
      return c.Status(500).SendString("token not valid")
    }

    adminIDBytes, _ := json.Marshal(claims.AdminID)
    adminIDJSON := string(adminIDBytes)
    c.Locals("adminID", adminIDJSON)
    
    schoolIDBytes, _ := json.Marshal(claims.SchoolID)
    schoolIDJSON := string(schoolIDBytes)
    c.Locals("schoolID", schoolIDJSON)

  }

  if (token == "") {
    return c.Status(500).SendString("no token")
  }

  return c.Next()
}

func postLogin(c *fiber.Ctx) error {

  // getting body and unmarshalling it into a body variable
  var body map[string]string
  json.Unmarshal(c.Body(), &body)

  // getting the db collection
  collection, err := db.GetCollection("admins")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // getting the admin
  var admin models.Admin
  if err = collection.FindOne(context.Background(), bson.M{"email": body["email"]}).Decode(&admin); err != nil {
    return c.Status(401).JSON(bson.M{
      "message": "Nu există niciun administrator cu email-ul introdus.",
    })
  }

  // // if the admin doesn't have a password
  // if admin.Password == "" {
  //   hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body["password"]), 10)
  //   if err != nil {
  //     return c.Status(500).SendString(fmt.Sprintf("%v", err))
  //   }

  //   var modifiedAdmin models.Admin
  //   collection.FindOneAndUpdate(context.Background(), bson.M{"cnp": body["cnp"]}, bson.D{
  //     {Key: "$set", Value: bson.D{{Key: "password",Value: string(hashedPassword)}}},
  //   }).Decode(&modifiedAdmin)

    
  //   // jwt
  //   tokenString, err := utils.AdminGenerateToken(modifiedAdmin.AdminID, modifiedAdmin.SchooldID)
  //   if err != nil {
  //     return c.Status(500).SendString(fmt.Sprintf("%v", err))
  //   }

  //   return c.JSON(bson.M{
  //     "adminID": modifiedAdmin.AdminID,
  //     "schoolID": modifiedAdmin.SchooldID,
  //     "firstName": modifiedAdmin.FirstName,
  //     "lastName": modifiedAdmin.LastName,
  //     "email": modifiedAdmin.Email,
  //     "password": modifiedAdmin.Password,
  //     "token": tokenString,
  //   })
  // } else {
    hashedPassword := admin.Password

    compareErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(body["password"]))

    tokenString, err := utils.AdminGenerateToken(admin.AdminID, admin.SchooldID)
    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    if compareErr == nil {
      return c.JSON(bson.M{
        "adminID": admin.AdminID,
        "schoolID": admin.SchooldID,
        "firstName": admin.FirstName,
        "lastName": admin.LastName,
        "email": admin.Email,
        "password": admin.Password,
        "token": tokenString,
      })
    } else {
      return c.Status(401).JSON(bson.M{
        "message": "Nu ați introdus parola validă.",
      })
    }
  // } 
}