package controllers

import (
	"fmt"
	"context"
	"log"
	"net/http"
	"time"

	"go-gin-mongo-jwt/models"
	"go-gin-mongo-jwt/helpers"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New();

func hashPassword(password string) string {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 15);

	if err != nil {
		log.Panic(err);
	}

	return string(bytes);
}

func verifyPassword(userPassword string, providedPassword string) bool {

	isValid := true;

	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword));

	if err != nil {
		fmt.Println(err);
		isValid = false;
	}

	return isValid;
}

func Signup(c *gin.Context) {

	fmt.Println("Signing up");

	var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second);

	var user models.User;


	// Binds json payload to model `json` tag
	err := c.BindJSON(&user); 
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid or empty fields supplied.",
		});
		return;
	}

	defer cancel();

	// Validate user inputs with model
	validationErr := validate.Struct(user);

	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()});
		return;
	}

	// Check for existing email and phone number in database

	var existingUser models.User;

	_ = models.UserCollection.FindOne(ctx, bson.M{"email": user.Email, "phone": user.Phone }).Decode(&existingUser);

	if existingUser.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or phone number already exists"});
		return;
	}

	user.Created_at = time.Now();
	user.Updated_at = time.Now();

	user.ID = primitive.NewObjectID();

	token, refreshToken, err := helpers.GenerateAllTokens(&user);

	if err != nil {
		log.Fatal(err);
	}

	user.Token = token;
	user.Refresh_token = refreshToken;

	password := hashPassword(user.Password);

	user.Password = password;

	result, err := models.UserCollection.InsertOne(ctx, user);

	if err != nil {
		fmt.Print(err);
		c.JSON(http.StatusBadRequest, gin.H{"error": "user was not created"});
		return;
	}

	defer cancel();

	c.JSON(http.StatusOK, result);

}

func Login(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second);

	var user models.User;

	var foundUser models.User;

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()});
	}

	_ = models.UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser);

	if foundUser.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"});
	}

	defer cancel();

	passwordIsValid := verifyPassword(user.Password, foundUser.Password);

	if !passwordIsValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password incorrect"});
		return;
	}

	// token, refreshToken, _ := helpers.GenerateAllTokens(&foundUser);	

	// helpers.UpdateAllTokens(token, refreshToken, foundUser.ID.Hex());

	c.JSON(http.StatusOK, foundUser);
}