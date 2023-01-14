package main

import (
	"fmt"
	"net/http"

	"go-gin-mongo-jwt/routes"
	"go-gin-mongo-jwt/configs"

	"github.com/gin-gonic/gin"
)

func main() {

	config := configs.Load();

	port := config.PORT;

	fmt.Print(port);

	// if port == "" {
	// 	port = "8000"
	// }

	r := gin.New();

	r.Use(gin.Logger());

	r.GET("/", func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Gin Mongo JWT server"});
	});

	r.GET("/api/v1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted"});
	});

	routes.AuthRoutes(r);
	routes.UserRoutes(r);

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Invalid route"});
	});

	r.Run(":" + port); 
}