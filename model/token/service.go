package token

import (
	"auth/model"
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

var DataSource *model.DataSource

func (refreshToken *RefreshToken) RefreshTokens(GUID string, ctx context.Context) (*TokenPair, error) {

	tokenPair, err := refreshToken.findOne(ctx)
	if err != nil {
		return nil, err
	}

	result, err := DataSource.AccountCollection.DeleteOne(ctx, bson.M{"_id": tokenPair.ID})
	if err != nil {
		return nil, err
	}

	if result.DeletedCount != 1 {
		return nil, err
	}

	tokenPair, err = GenerateTokens(GUID, ctx)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (refreshToken RefreshToken) Delete(ctx context.Context) error {

	tokenPair, err := refreshToken.findOne(ctx)
	if err != nil {
		return err
	}

	tokenPair.RefreshToken = ""
	_, err = DataSource.AccountCollection.ReplaceOne(
		ctx,
		bson.M{"_id": tokenPair.ID},
		tokenPair)

	if err != nil {
		return err
	}

	return nil
}

func (refreshToken RefreshToken) findOne(ctx context.Context) (*TokenPair, error) {
	guid, err := GetSubClaims(string(refreshToken))
	if err != nil {
		return nil, err
	}

	cursor, err := DataSource.AccountCollection.Find(ctx, bson.M{"account_guid": guid})
	if err != nil {
		return nil, err
	}

	var tokenPairs []TokenPair
	if err = cursor.All(context.TODO(), &tokenPairs); err != nil {
		return nil, err
	}

	for _, tokenPair := range tokenPairs {

		hashedRefreshToken := tokenPair.RefreshToken

		err = bcrypt.CompareHashAndPassword([]byte(hashedRefreshToken), []byte(string(refreshToken)))
		if err == nil {
			return &tokenPair, nil
		}
	}

	return nil, errors.New("cat'n perfrom token from database")
}

func (refreshToken RefreshToken) hashAndSalt() (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func DeleteAllRefreshTokens(GUID string, ctx context.Context) error {

	_, err := DataSource.AccountCollection.UpdateMany(
		ctx,
		bson.M{"account_guid": GUID},
		bson.D{
			{"$set", bson.D{{"refresh_token", ""}}},
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func GenerateTokens(GUID string, ctx context.Context) (*TokenPair, error) {

	accessToken := jwt.New(jwt.GetSigningMethod("HS512"))
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["sub"] = GUID
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	// adding claims

	at, err := accessToken.SignedString([]byte(os.Getenv("token_secret")))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(jwt.GetSigningMethod("HS512"))
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = GUID
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	// adding claims

	rt, err := refreshToken.SignedString([]byte(os.Getenv("token_secret")))
	if err != nil {
		return nil, err
	}

	tokens := TokenPair{GUID: GUID, AccessToken: at, RefreshToken: rt}
	err = tokens.save(ctx)
	if err != nil {
		return nil, err
	}

	return &tokens, nil
}

func (tokenPair TokenPair) save(ctx context.Context) error {

	hashedRefreshToken, err := RefreshToken(tokenPair.RefreshToken).hashAndSalt()
	if err != nil {
		return err
	}

	tokenPairToSave := bson.D{
		{"account_guid", tokenPair.GUID},
		{"access_token", tokenPair.AccessToken},
		{"refresh_token", hashedRefreshToken}}
	_, err = DataSource.AccountCollection.InsertOne(ctx, tokenPairToSave)
	if err != nil {
		return err
	}

	return nil
}

func GetSubClaims(tokenStr string) (string, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_secret")), nil
	})
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)

	guid, ok := claims["sub"]
	if !ok {
		log.Print("missing sub in claims")

		return "", errors.New("missing sub in claims")
	}

	return guid.(string), nil
}
