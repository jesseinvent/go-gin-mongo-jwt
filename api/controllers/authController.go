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

func hashPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 15);
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func verifyPassword(userPassword string, providedPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword));

	//return nil if passwords are thesame
	return err == nil;
}

func Signup(c *gin.Context) {

	fmt.Println("Signing up")

	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
	defer cancel()

	var user models.User


	// Binds json payload to model `json` tag
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid or empty fields supplied.",
		});
		return
	}

	

	// Validate user inputs with model
	if validationErr := validate.Struct(user); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()});
		return
	}

	// Check for existing email and phone number in database

	var existingUser models.User;

	_ = models.UserCollection.FindOne(ctx, bson.M{"email": user.Email, "phone": user.Phone }).Decode(&existingUser);

	if existingUser.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or phone number already exists"});
		return
	}

	user.Created_at = time.Now();
	user.Updated_at = time.Now();

	user.ID = primitive.NewObjectID();

	token, refreshToken, err := helpers.GenerateAllTokens(&user);
	if err != nil {
		fmt.Println(err.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an error occcured please try again"});
		return
	}

	user.Token = token
	user.Refresh_token = refreshToken

	password, err := hashPassword(user.Password)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password format"});
		return
	}
	
	user.Password = password

	result, err := models.UserCollection.InsertOne(ctx, user);
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "user was not created"});
		return
	}


	c.JSON(http.StatusOK, result);
}

func Login(c *gin.Context) {

	fmt.Println("Login")

	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
	defer cancel()

	var user, foundUser models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ = models.UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser);

	if foundUser.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"});
		return
	}


	if passwordIsValid := verifyPassword(user.Password, foundUser.Password); !passwordIsValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password incorrect"});
		return
	}

	// token, refreshToken, _ := helpers.GenerateAllTokens(&foundUser);	

	// helpers.UpdateAllTokens(token, refreshToken, foundUser.ID.Hex());

	c.JSON(http.StatusOK, foundUser);
}
