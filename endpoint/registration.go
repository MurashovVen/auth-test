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

var Registration = func(w http.ResponseWriter, r *http.Request) {

	log.Print("got [/api/account/register] request")

	account := &a.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		log.Print(err)

		u.Respond(w, u.Message(false, "invalid request"))
		return
	}

	log.Printf("registering user.guid: %s", account.GUID)
	err = account.Register(context.TODO())
	if err != nil {
		log.Print(err)

		w.WriteHeader(http.StatusForbidden)
		u.Respond(w, u.Message(false, "can't register user"))
		return
	}

	log.Printf("generating tokens for user.guid: %s", account.GUID)
	tokens, err := t.GenerateTokens(account.GUID, context.TODO())
	if err != nil {
		log.Print(err)

		w.WriteHeader(http.StatusInternalServerError)
		u.Respond(w, u.Message(false, "can't generate token"))
		return
	}

	u.Respond(w, tokens)
}
