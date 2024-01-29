package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"task/api/database"
	"task/api/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authenticate(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Check exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// Find the user with the same token
		var user models.User
		database.DB.First(&user, "id =?", claims["sub"])
		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// Attach to request
		c.Set("user", user)
		// Continue
		c.Next()
		fmt.Println(claims["sub"], claims["exp"])
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	fmt.Println("Authenticating user")
}
func GetLoginInfo(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	fmt.Println(err)
	if err == nil {
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			log.Fatal(err)
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Check exp
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.AbortWithStatus(http.StatusUnauthorized)
			}
			// Find the user with the same token
			var user models.User
			database.DB.First(&user, "id =?", claims["sub"])
			if user.ID == 0 {
				c.AbortWithStatus(http.StatusUnauthorized)
			}
			// Attach to request
			c.Set("user", user)
			// Continue
			c.Next()
			fmt.Println(claims["sub"], claims["exp"])
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)

		}
	} else {
		var user = models.User{}
		c.Set("user", user)
		c.Next()
	}
	fmt.Println("Authenticating user")
}
