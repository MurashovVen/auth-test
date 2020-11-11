package endpoint

import (
	a "auth/model/account"
	t "auth/model/token"
	u "auth/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

var AccessRefreshTokens = func(w http.ResponseWriter, r *http.Request) {

	log.Print("got [/api/account/login] request")

	accountReq := &a.Account{}
	err := json.NewDecoder(r.Body).Decode(accountReq)
	if err != nil {
		log.Print(err)

		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(false, "invalid request"))
		return
	}

	log.Print("loading account from database")
	accountDb, err := a.LoadByUsername(accountReq.Username, context.TODO())
	if err != nil {
		log.Print(err)

		w.WriteHeader(http.StatusNotFound)
		u.Respond(w, u.Message(false, "can't perform account from database"))
		return
	}

	log.Print("verification credentials")
	if accountDb.Password != accountReq.Password || accountDb.GUID != accountReq.GUID {
		log.Print(err)

		w.WriteHeader(http.StatusForbidden)
		u.Respond(w, u.Message(false, "invalid auth"))
		return
	}

	log.Printf("generating tokens for user.guid: %s", accountReq.GUID)
	tokens, err := t.GenerateTokens(accountDb.GUID, context.TODO())
	if err != nil {
		log.Print(err)

		w.WriteHeader(http.StatusInternalServerError)
		u.Respond(w, u.Message(false, "can't perform tokens"))
		return
	}

	u.Respond(w, tokens)
}
