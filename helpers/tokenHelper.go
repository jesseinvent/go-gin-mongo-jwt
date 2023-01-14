package helpers

import(
	"context"
	"fmt"
	"log"
	"time"

	"go-gin-mongo-jwt/models"
	"go-gin-mongo-jwt/configs"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/golang-jwt/jwt/v4"
)

type SignedDetails struct {
	Email 			string
	First_name 		string
	Last_name 		string
	ID 				string
	User_type		string
	jwt.RegisteredClaims
}

var config = configs.Load();

var SECRET_KEY = config.SECRET;

func GenerateAllTokens(user *models.User) (string, string, error) {

	fmt.Println("Generating tokens");

	claims := &SignedDetails{
		Email: user.Email,
		First_name: user.First_name,
		Last_name: user.Last_name,
		ID: user.ID.Hex(),
		User_type: user.User_type,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{time.Now().Add(time.Hour * 24)},
		},
	}

	refreshClaims := SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{time.Now().Add(time.Hour * 24)},
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY));

	if err != nil {
		err = fmt.Errorf("failed creating token: %s", err);
		return "", "", err;
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY));

	if err != nil {
		err = fmt.Errorf("failed creating refresh token: %s", err);
		return "", "", err;
	}

	return token, refreshToken, nil;
} 

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	// Claims: Encoded information in the token

	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{},
	func(token *jwt.Token)(interface{}, error){
		return []byte(SECRET_KEY), nil
	});

	if err != nil {
		msg = fmt.Sprintf("error decoding token: %s", err.Error());
		return;
	}

	claims, ok := token.Claims.(*SignedDetails);

	if !ok {
		msg = fmt.Sprintf("the token is invalid: %s", err.Error());
		return;
	}

	tokenExpiry := claims.ExpiresAt.Unix();

	if tokenExpiry < time.Now().Unix() {
		msg = fmt.Sprintf("token is expired: %s", err.Error());
		return;
	}

	return claims, msg;
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	ctx, cancle := context.WithTimeout(context.Background(), 100 * time.Second);

	var updateObj primitive.D;

	updateObj = append(updateObj, bson.E{"token", signedToken});
	
	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken});

	Updated_at := time.Now();

	updateObj = append(updateObj, bson.E{"updated_at", Updated_at});

	upsert := true;

	filter := bson.M{"user_id": userId};

	opt := &options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := models.UserCollection.UpdateOne(ctx, filter, bson.D{
		{"$set", updateObj},
	}, opt);

	defer cancle();

	if err != nil {
		log.Fatal(err);
	}
}