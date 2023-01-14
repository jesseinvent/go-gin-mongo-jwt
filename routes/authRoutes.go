package routes

import(
	"go-gin-mongo-jwt/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	routes := r.Group("/api/v1/auth");

	routes.POST("/signup", controllers.Signup);
	routes.POST("/login", controllers.Login);
}