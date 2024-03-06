package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserInfo struct {
	ID          string `json:"id"`
	SessionGUID string `json:"session"`
}

type JwtUserClaims struct {
	jwt.RegisteredClaims
	User *UserInfo
}

// NewRefreshToken generates a new refresh token using the provided secret, expiration time (in minutes), and user information.
// It returns the refresh token as a string.
func NewRefreshToken(secret string, userInfo *UserInfo) (string, error) {
	claims := JwtUserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationTime) * time.Minute)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
		User: userInfo,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString([]byte(secret))
	return refreshToken, err
}

// NewAccessToken generates a new access token with the given parameters.
// It takes a secret string, expiration time in minutes, and an ID string.
// Returns the generated access token as a string and any error encountered.
func NewAccessToken(secret string, id string) (string, error) {
	claims := JwtUserClaims{
		User: &UserInfo{
			ID: id,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: id,
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationTime) * time.Minute)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return accessToken.SignedString([]byte(secret))
}

func VerifyToken(secret string, token string) (*UserInfo, bool) {
	t, err := jwt.ParseWithClaims(token, &JwtUserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse token", err)
		return nil, false
	}

	// expTime, err := t.Claims.GetExpirationTime()
	// if err != nil {
	// 	return nil, false
	// }

	// if !t.Valid || expTime.Before(time.Now()) {
	// 	fmt.Println("expires token on invalid token")
	// 	return nil, false
	// }

	claims := t.Claims.(*JwtUserClaims)
	if claims == nil {
		return nil, false
	}
	return claims.User, true
}

func ParseRefreshToken(secret, refresh string) (*UserInfo, error) {
	t, err := jwt.ParseWithClaims(refresh, &JwtUserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(*JwtUserClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse refresh token")
	}

	return claims.User, nil
}

type TokensPair struct {
	AccessToken  string
	RefreshToken string
}

func GetAuthTokenFromBearerToken(token string) (string, error) {
	if len(token) < 7 || token[:7] != "Bearer " {
		return "", fmt.Errorf("invalid token format")
	}
	return token[7:], nil
}
