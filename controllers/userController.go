package controllers

import (
	_"fmt"
	"log"
	"strconv"
	"net/http"
	"time"
	"context"

	"go-gin-mongo-jwt/models"
	"go-gin-mongo-jwt/helpers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

func GetUser(c *gin.Context) {
	userId := c.Param("user_id");
	id := c.GetString("ID");

	if userId != id {
		c.JSON(http.StatusBadRequest, gin.H{"err": "unauthorized to access resource"});
	}

	if err := helpers.CheckIfUserIsAdmin(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()});
		return;
	} 

	var ctx, cancle = context.WithTimeout(context.Background(), 100 * time.Second);

	var user models.User;

	objId, _ := primitive.ObjectIDFromHex(userId);

	err := models.UserCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user); 

	defer cancle();

	// No user found
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()});
		return;
	}

	c.JSON(http.StatusOK, user);

	return;
}

func GetUsers(c *gin.Context) {

	// Check if user is an admin
	if err := helpers.CheckIfUserIsAdmin(c); err != nil {
		
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()});

		return;
	} 

	var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second);

	recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"));

	if err != nil || recordPerPage < 1{
		recordPerPage = 10;
	}

	// Pagination
	page, err := strconv.Atoi(c.Query("page"));

	if err != nil || page < 1 {
		page = 1;
	}

	startIndex :=  (page - 1) * recordPerPage;

	// MongoDB aggregation pipeline
	matchStage := bson.D{{"$match", bson.D{{}}}};

	// _id signals where the data will be grouped on
	groupStage := bson.D{
			{
				"$group", bson.D{
					{"_id", bson.D{{"_id", "null"}}},
					{"total_count", bson.D{{"$sum", 1}}}, 
					{"data", bson.D{{"$push", "$$ROOT"}}},
			},
		},
	};

	projectStage := bson.D{
		{
			"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{
					"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex,recordPerPage}}},
				},
			},
		},
	};

	result, err := models.UserCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage, groupStage, projectStage,
	});

	defer cancel();

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error when listing user records"});
	}

	var allUsers []bson.M;

	err = result.All(ctx, &allUsers);

	if err != nil {
		log.Fatal(err);
	}

	c.JSON(http.StatusOK, allUsers[0]);

} 