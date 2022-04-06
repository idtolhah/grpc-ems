package auth

import (
	"bff/client"
	"strings"

	"github.com/gin-gonic/gin"
)

// IsAuthenticated validates token and authorizes users
func IsAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(403, "No Authorization header provided")
			c.Abort()
			return
		}

		extractedToken := strings.Split(clientToken, "Bearer ")

		if len(extractedToken) == 2 {
			clientToken = strings.TrimSpace(extractedToken[1])
		} else {
			c.JSON(400, "Incorrect Format of Authorization Token")
			c.Abort()
			return
		}

		jwtWrapper := JwtWrapper{
			SecretKey: client.GetEnv("APP_JWT_SECRET"),
		}

		claims, err := jwtWrapper.ValidateToken(clientToken)
		if err != nil {
			c.JSON(401, err.Error())
			c.Abort()
			return
		}

		c.Set("UserId", claims.Id)

		c.Next()

	}
}
