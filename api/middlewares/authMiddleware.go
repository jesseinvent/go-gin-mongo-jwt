package middlewares

import(
	_ "fmt"
	"net/http"
	"strings"

	"go-gin-mongo-jwt/helpers"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.Request.Header.Get("Authorization");

		if authHeader == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No Authorization header provided." });
			c.Abort(); // Stops request
			return;
		}

		token := strings.Split(authHeader, " ")[1];

		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No Auth token provided in header." });
			c.Abort(); // Stops request
			return;
		}

		claims, err := helpers.ValidateToken(token);

		if err != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err});
			c.Abort();
			return;
		}
	
		c.Set("email", claims.Email);
		c.Set("first_name", claims.First_name);
		c.Set("last_name", claims.Last_name);
		c.Set("ID", claims.ID);
		c.Set("user_type", claims.User_type);
		c.Next();
	}
}

