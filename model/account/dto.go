package account

type Account struct {
	GUID     string `bson:"_id"       json:"GUID"`
	Username string `bson:"username"  json:"username"`
	Password string `bson:"password"  json:"password"`
}
