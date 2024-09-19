package tokens

import (
	"log"
	"os"
	"time"

	"github.com/adil-mubarak/ecommerce/database"
	"github.com/adil-mubarak/ecommerce/models"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type SignedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.StandardClaims
}

var SECRET_KEY = os.Getenv("SECRET_LOVE")
var DB *gorm.DB = database.DB

// TokenGenerator generates the JWT token and refresh token
func TokenGenerator(email, firstname, lastname, uid string) (signedToken, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_Name: firstname,
		Last_Name:  lastname,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 24).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 168).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

// ValidateToken validates the JWT token
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "The token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token is expired"
		return
	}

	return claims, ""
}

// UpdateAllTokens updates the tokens in the database for a user
func UpdateAllTokens(signedToken, signedRefreshToken, userID string) {
	var user models.User
	result := DB.Where("uid = ?", userID).First(&user)
	if result.Error != nil {
		log.Println(result.Error)
		return
	}

	user.Token = signedToken
	user.RefreshToken = signedRefreshToken
	user.UpdatedAt = time.Now()

	result = DB.Save(&user)
	if result.Error != nil {
		log.Println(result.Error)
	}
}
