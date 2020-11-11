package endpoint

import (
	a "auth/app"
	t "auth/model/token"
	u "auth/utils"
	"context"
	"log"
	"net/http"
	"strings"
)

var Refresh = func(w http.ResponseWriter, r *http.Request) {

	tokenHeader := r.Header.Get("refresh_token")

	log.Printf("got refresh [/api/account/refresh] request with refresh token : %s", tokenHeader)

	ok, response := a.Validation(tokenHeader)
	if !ok {
		log.Print("invalid refresh token")

		w.WriteHeader(http.StatusForbidden)
		u.Respond(w, response)
		return
	}

	refreshToken := t.RefreshToken(tokenHeader)

	log.Print("extracting claims from token")
	guid, err := t.GetSubClaims(tokenHeader)
	if err != nil {
		log.Print(err)

		w.WriteHeader(http.StatusForbidden)
		u.Respond(w, u.Message(false, "invalid refresh token"))
		return
	}

	log.Print("refreshing token pair")
	tokens, err := refreshToken.RefreshTokens(guid, context.TODO())
	if err != nil {
		log.Print(err)

		w.WriteHeader(http.StatusForbidden)
		u.Respond(w, u.Message(false, "can`t refresh tokens"))
		return
	}

	u.Respond(w, tokens)
}

var DelRefreshToken = func(w http.ResponseWriter, r *http.Request) {

	response := make(map[string]interface{})
	tokenHeader := r.Header.Get("refresh_token")

	log.Printf("got delete [/api/account/refresh] request with refresh token : %s", tokenHeader)

	ok, response := a.Validation(tokenHeader)
	if !ok {
		log.Print("invalid token")

		w.WriteHeader(http.StatusForbidden)
		u.Respond(w, response)
		return
	}

	refreshToken := t.RefreshToken(tokenHeader)

	log.Printf("deleting refresh token: %s", tokenHeader)
	err := refreshToken.Delete(context.TODO())
	if err != nil {
		log.Print(err)

		w.WriteHeader(http.StatusInternalServerError)
		u.Respond(w, u.Message(false, "refresh token doesn't exist"))
		return
	}

	u.Respond(w, u.Message(true, "token deleted"))
}

var DelAllRefreshTokens = func(w http.ResponseWriter, r *http.Request) {

	log.Printf("got delete [/api/account/refresh/all] request with refresh token")

	tokenStr := strings.Split(r.Header.Get("Authorization"), " ")[1]

	log.Print("extracting claims from token")
	guid, err := t.GetSubClaims(tokenStr)
	if err != nil {
		log.Print(err)

		w.WriteHeader(http.StatusForbidden)
		u.Respond(w, u.Message(true, "can't extract sub from claims"))
		return
	}

	log.Printf("deleting all refesh tokens from user.guid: %s", guid)
	err = t.DeleteAllRefreshTokens(guid, context.TODO())
	if err != nil {
		log.Print(err)

		w.WriteHeader(http.StatusInternalServerError)
		u.Respond(w, u.Message(true, "can't delete tokens from database"))
		return
	}

	u.Respond(w, u.Message(true, "tokens deleted"))
}
