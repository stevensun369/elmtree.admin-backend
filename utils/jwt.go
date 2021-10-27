package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWTKey = []byte("123456")

type AdminClaims struct {
  AdminID string `json:"adminID"`
  SchoolID string `json:"schoolID"`
  jwt.StandardClaims
}

func AdminGenerateToken(id string, schoolID string) (tokenString string, err error) {
  // one year has 8760 hours
  expirationTime := time.Now().Add(8760 * time.Hour)

  claims := &AdminClaims {
    AdminID: id,
    SchoolID: schoolID,
    StandardClaims: jwt.StandardClaims {
      ExpiresAt: expirationTime.Unix(),
    },
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  tokenString, err = token.SignedString(JWTKey)

  return tokenString, err
}