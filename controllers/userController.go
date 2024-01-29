package controllers

import (
	"fmt"
	"net/http"
	"os"
	"task/api/database"
	"task/api/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func IndexGET(c *gin.Context) {
	user, _ := c.Get("user")
	empty := models.User{}
	if user == empty {
		c.HTML(http.StatusOK, "index.html", nil)
		return
	}
	c.HTML(http.StatusOK, "index.html", user)

}
func SignUp(c *gin.Context) {
	if c.Request.Method == "POST" {
		var body struct {
			Username string `json:"username" form:"username"`
			Email    string `json:"email" form:"email"`
			Password string `json:"password" form:"password"`
		}
		if c.ShouldBind(&body) != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "Failed to read body",
			})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to hash password"})
			return
		}

		user := models.User{
			Username:  body.Username,
			Email:     body.Email,
			Password:  string(hash),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now()}
		result := database.DB.Create(&user)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": result.Error})
			return
		}
		msg := "User created successfully"
		c.JSON(http.StatusOK, gin.H{"message": msg})
	} else {
		c.HTML(http.StatusOK, "userRegister.html", gin.H{})
	}

}
func SignUpGET(c *gin.Context) {
	user, _ := c.Get("user")
	empty := models.User{}
	if user == empty {
		c.HTML(http.StatusOK, "userRegister.html", nil)
		return
	}
	c.HTML(http.StatusOK, "userRegister.html", user)

}

func Login(c *gin.Context) {
	if c.Request.Method == "POST" {
		var body struct {
			Email    string `json:"email" form:"email"`
			Password string `json:"password" form:"password"`
		}
		if c.ShouldBind(&body) != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "Failed to read body",
			})
			return
		}

		// Look up requested user
		var user models.User
		database.DB.First(&user, "email = ?", body.Email)
		if user.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "User not found with this Email:" + " " + body.Email})
			return
		}
		fmt.Println("This is user :\n", user.Username, " \n", user.Email, " \n", user.Password, " \nEND")
		// Compare hash with saved user pass
		var err error
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Error comparing passwords" + " " + "Error:" + err.Error()})
			return
		}
		// Generate a jwt token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		})
		// Sign and get the complete encoded token as a string using the secret

		tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to create token " + " " + "Error:" + err.Error()})
			return
		}
		// Save to cookie
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
		c.HTML(http.StatusOK, "userLogin.html", body)

	} else {
		c.HTML(http.StatusOK, "userLogin.html", gin.H{})
	}
}
func LoginGET(c *gin.Context) {
	user, _ := c.Get("user")
	empty := models.User{}
	if user == empty {
		c.HTML(http.StatusOK, "userLogin.html", nil)
		return
	}
	c.HTML(http.StatusOK, "userLogin.html", user)

}
func UserProfileGET(c *gin.Context) {
	userId := c.Param("userId")
	user, _ := c.Get("user")
	empty := models.User{}
	if user == empty {
		c.HTML(http.StatusOK, "user.html", nil)
		return
	} else {
		var user models.User
		database.DB.First(&user, "id =?", userId)
		c.HTML(http.StatusOK, "user.html", user)
	}

}
func UpdateUser(c *gin.Context) {
	if c.Request.Method == "POST" && c.PostForm("_method") == "PUT" {
		user, _ := c.Get("user")
		empty := models.User{}
		if user == empty {
			c.HTML(http.StatusOK, "user.html", nil)
			return
		}
		fmt.Println(c.Request.Method)
		if c.Request.Method == "POST" {
			var body struct {
				Username string `json:"username" form:"username"`
				Email    string `json:"email" form:"email"`
				Password string `json:"password" form:"password"`
			}

			if c.ShouldBind(&body) != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"Error": "Failed to read body",
				})
				return
			}

			hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to hash password"})
				return
			}

			user, _ := c.Get("user")
			database.DB.Model(&user).Where("id = ?", user.(models.User).ID).Select("username", "email", "password", "updated_at").Updates(models.User{Username: body.Username, Email: body.Email, Password: string(hash), UpdatedAt: time.Now()})
			// database.DB.Model(&user).Updates(models.User{Username: body.Username, Email: body.Email, Password: body.Password})
			msg := "User" + user.(models.User).Username + " updated successfully"
			c.JSON(http.StatusOK, gin.H{"message": msg})

		} else {
			c.HTML(http.StatusOK, "user.html", user)
		}
	} else if c.Request.Method == "POST" && c.PostForm("_method") == "DELETE" {

		user, _ := c.Get("user")
		database.DB.Unscoped().Model(&user).Where("id = ?", user.(models.User).ID).Delete(&models.User{})
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request"})
	}
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"message": user})

}
