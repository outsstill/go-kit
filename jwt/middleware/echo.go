package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	myjwt "github.com/outsstill/go-kit/jwt"
)

func EchoJWT(j *myjwt.JWT) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {

			auth := c.Request().
				Header.Get("Authorization")

			if auth == "" {
				return c.JSON(
					http.StatusUnauthorized,
					map[string]any{
						"msg": "token required",
					},
				)
			}

			parts := strings.SplitN(auth, " ", 2)

			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(
					http.StatusUnauthorized,
					map[string]any{
						"msg": "invalid authorization",
					},
				)
			}

			claims, err := j.ParserToken(parts[1])
			if err != nil {
				return c.JSON(
					http.StatusUnauthorized,
					map[string]any{
						"msg": err.Error(),
					},
				)
			}

			c.Set("claims", claims)

			return next(c)
		}
	}
}
