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

	// mongo
	"go.mongodb.org/mongo-driver/bson"

	// bcrypt
	"golang.org/x/crypto/bcrypt"
)

func AdminMiddleware(c *fiber.Ctx) error {
  var token string

  authHeader := c.Get("Authorization")

  if string(authHeader) != "" && strings.HasPrefix(string(authHeader), "Bearer") {
    token = strings.Fields(string(authHeader))[1]

    // we're just parsing the token: maybe I will put it in the utils
    claims := &utils.AdminClaims{}
    tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface {}, error) {
      return utils.JWTKey, nil
    })

    utils.CheckError(c, err)

    if !tkn.Valid {
      return utils.MessageError(c, "token not valid")
    }

    utils.SetLocals(c, "adminID", claims.AdminID)
    utils.SetLocals(c, "schoolID", claims.SchoolID)

  }

  if (token == "") {
    return utils.MessageError(c, "no token")
  }

  return c.Next()
}

// @desc   login
// @route  POST /api/admin/login
// @access Private
func postLogin(c *fiber.Ctx) error {

  // getting body and unmarshalling it into a body variable
  var body map[string]string
  json.Unmarshal(c.Body(), &body)

  // getting the admin
  var admin models.Admin
  if err := db.Admins.FindOne(context.Background(), bson.M{"email": body["email"]}).Decode(&admin); err != nil {
    return c.Status(401).JSON(bson.M{
      "message": "Nu există niciun administrator cu email-ul introdus.",
    })
  }

  hashedPassword := admin.Password

  compareErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(body["password"]))

  tokenString, err := utils.AdminGenerateToken(admin.AdminID, admin.SchooldID)
  utils.CheckError(c, err)

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
    return utils.MessageError(c, "Nu ați introdus parola validă.")
  }
}