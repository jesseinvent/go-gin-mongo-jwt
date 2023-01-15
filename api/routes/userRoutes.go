package routes

import(
	"go-gin-mongo-jwt/controllers"
	"go-gin-mongo-jwt/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine){

	routes := r.Group("/api/v1/users");

	routes.Use(middlewares.Authenticate());

	routes.GET("/", controllers.GetUsers);
	routes.GET("/:user_id", controllers.GetUser); 
}