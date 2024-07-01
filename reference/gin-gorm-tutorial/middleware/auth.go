package middleware

import (
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)


func IsAuth() gin.HandlerFunc {
	return checkJWT(false)
}
func IsAdmin() gin.HandlerFunc {
	return checkJWT(true)
}

func checkJWT(middlewareAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader," ")

		if len(bearerToken) == 2 {
			token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("SUPER_SECRET")), nil
			})

			if err == nil && token.Valid {
				claims := token.Claims.(jwt.MapClaims)
				userRole := bool(claims["user_role"].(bool))

				c.Set("jwt_user_id", claims["user_id"])
				c.Set("jwt_isAdmin", claims["user_role"])

				if middlewareAdmin == true && userRole == false {

					c.JSON(422, gin.H{
						"msg": "admin only allowed",
						"error": err,
					})
					c.Abort()
					return
				}

			} else {
				c.JSON(422, gin.H{
					"msg": "error",
					"error": err,
				})
				c.Abort()
				return

			}
		} else {

				c.JSON(422, gin.H{
					"msg": "Authorization token is not provided",
				})
				c.Abort()
				return
		}
		
	}
	
}