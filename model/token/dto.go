package token

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenPair struct {
	ID           primitive.ObjectID `bson:"_id"            json:"-"`
	GUID         string             `bson:"account_guid"   json:"GUID"`
	AccessToken  string             `bson:"access_token"   json:"access_token"`
	RefreshToken string             `bson:"refresh_token"  json:"refresh_token"`
}

type AccessToken string
type RefreshToken string
