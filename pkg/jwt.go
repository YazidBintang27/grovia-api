package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type JWTClaim struct {
	UserID     int    `json:"user_id"`
	LocationID int    `json:"location_id"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID, locationID int, role string) (accessToken, refreshToken string, err error) {
	accessClaims := &JWTClaim{
		UserID:     userID,
		Role:       role,
		LocationID: locationID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "grovia-api",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = at.SignedString(jwtKey)

	if err != nil {
		return
	}

	refreshClaims := &JWTClaim{
		UserID:     userID,
		Role:       role,
		LocationID: locationID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			Issuer:    "grovia-api",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = rt.SignedString(jwtKey)

	if err != nil {
		return
	}

	return accessToken, refreshToken, nil
}

func ValidateToken(signedToken string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return jwtKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
