package auth

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// TokenParser represents JWT token parser
type TokenParser interface {
	ParseToken(string) (*jwt.Token, error)
}

// Middleware makes JWT implement the Middleware interface.
func Middleware(tokenParser TokenParser) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := tokenParser.ParseToken(c.Request().Header.Get("Authorization"))
			if err != nil || !token.Valid {
				return c.NoContent(http.StatusUnauthorized)
			}

			claims := token.Claims.(jwt.MapClaims)

			userName := claims["userName"].(string)
			userEmail := claims["userEmail"].(string)

			c.Set("name", userName)
			c.Set("user_email", userEmail)

			return next(c)
		}
	}
}
