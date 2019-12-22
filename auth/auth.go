package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/emanpicar/minimart-api/settings"
	"github.com/gorilla/context"

	jwt "github.com/dgrijalva/jwt-go"
)

type (
	Manager interface {
		Authenticate(body io.ReadCloser) (string, error)
		ValidateRequest(r *http.Request) error
	}

	authHandler struct{}

	User struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
)

func NewManager() Manager {
	return &authHandler{}
}

func (a *authHandler) Authenticate(body io.ReadCloser) (string, error) {
	var user User
	if err := json.NewDecoder(body).Decode(&user); err != nil {
		return "", err
	}

	// TODO validate user and pass in DB
	// code here <--

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
	})

	tokenString, err := token.SignedString([]byte(settings.GetTokenSecret()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *authHandler) ValidateRequest(r *http.Request) error {
	authorizationHeader := r.Header.Get("authorization")
	if authorizationHeader == "" {
		return errors.New("An authorization header is required")
	}

	bearerToken := strings.Split(authorizationHeader, " ")
	if len(bearerToken) != 2 {
		return errors.New("Cannot parse authorization header")
	}

	token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Cannot parse authorization header")
		}
		return []byte(settings.GetTokenSecret()), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("Invalid authorization token")
	}

	context.Set(r, "tokenClaims", token.Claims)

	return nil
}
