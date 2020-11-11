package account

import (
	"auth/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

var DataSource *model.DataSource

func (account *Account) Register(ctx context.Context) error {

	_, err := DataSource.AccountCollection.InsertOne(ctx, account)
	if err != nil {
		return err
	}

	return nil
}

func LoadByUsername(username string, ctx context.Context) (*Account, error) {

	var account Account
	err := DataSource.AccountCollection.FindOne(ctx, bson.M{"username": username}).Decode(&account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}
