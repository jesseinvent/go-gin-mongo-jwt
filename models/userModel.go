package models

import(
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go-gin-mongo-jwt/database"
)

type User struct {
	ID				primitive.ObjectID	`bson:"_id" json:"id"`
	First_name		string				`bson:"first_name" json:"first_name" validate:"required,min=2,max=100"`
	Last_name		string				`bson:"last_name" json:"last_name" validate:"required,min=2,max=100"`
	Password		string				`bson:"password" json:"password" validate:"required,min=6"`
	Email			string				`bson:"email" json:"email" validate:"email,required" `
	Phone			string				`bson:"phone" json:"phone" validate:"required,max=100"`
	Token			string				`bson:"token" json:"token"`
	User_type		string				`bson:"user_type" json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Refresh_token	string				`bson:"refresh_token" json:"refresh_token"`
	Created_at		time.Time			`bson:"created_at" json:"created_at"`
	Updated_at		time.Time			`bson:"updated_at" json:"updated_at"`
}


var UserCollection *mongo.Collection = database.OpenCollection(database.Client, "user");