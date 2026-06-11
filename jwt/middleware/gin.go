package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	myjwt "github.com/outsstill/go-kit/jwt"
)

func GinJWT(j *myjwt.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {

		auth := c.GetHeader("Authorization")

		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "token required",
			})
			return
		}

		parts := strings.SplitN(auth, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "invalid authorization",
			})
			return
		}

		claims, err := j.ParserToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": err.Error(),
			})
			return
		}

		c.Set("claims", claims)

		c.Next()
	}
}
