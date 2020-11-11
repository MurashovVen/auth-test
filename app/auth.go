package app

import (
	u "auth/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
	"strings"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/account/register", "/api/account/login", "/api/account/refresh"}
		requestPath := r.URL.Path

		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		log.Print("authenticate connection")

		tokenHeader := r.Header.Get("Authorization")
		split := strings.Split(tokenHeader, " ")
		if len(split) != 2 {
			log.Print("invalid/malformed auth token")

			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, u.Message(false, "invalid/malformed auth token"))
			return
		}

		tokenStr := split[1]
		ok, response := Validation(tokenStr)
		if !ok {
			log.Print("invalid auth token")

			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, response)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Validation(tokenStr string) (bool, map[string]interface{}) {

	if tokenStr == "" {
		return false, u.Message(false, "Missing auth token")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("token_secret")), nil
	})

	if err != nil {
		return false, u.Message(false, "Malformed authentication token")
	}

	if !token.Valid {
		return false, u.Message(false, "Token is not valid.")
	}

	return true, nil
}
