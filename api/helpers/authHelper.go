package helpers

import (
	_"fmt"
	"errors"
	"github.com/gin-gonic/gin"
)


func CheckIfUserIsAdmin(c *gin.Context) error {

	userType := c.GetString("user_type");

	if userType != "ADMIN" {
		err := errors.New("unauthorized to access this resource");

		return err;
	}

	return nil;
}
