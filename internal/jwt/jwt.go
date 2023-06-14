package jwt

import (
	"os"
	"time"

	"github.com/RamazanZholdas/KeyboardistSV2/utils"
	"github.com/dgrijalva/jwt-go/v4"
)

func NewJwtTokenWithClaims(email string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    email,
		ExpiresAt: &jwt.Time{Time: time.Now().Add(time.Hour * 24)},
	})

	token, err := claims.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		utils.LogError("Error signing token: ", err)
		return "", err
	}

	return token, nil
}

func ExtractTokenClaimsFromCookie(cookie string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		utils.LogError("Error parsing token: ", err)
		return nil, err
	}

	claims := token.Claims.(*jwt.StandardClaims)
	return claims, nil
}
