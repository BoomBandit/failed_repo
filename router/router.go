package router

import (
	// "net/http"
	"task/api/controllers"
	"task/api/middleware"

	// "task/api/models"

	"github.com/gin-gonic/gin"
)

func StartRouter() {
	// var pictures []models.Picture

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	// Index
	router.GET("/", middleware.GetLoginInfo, controllers.IndexGET)

	//Register routes
	router.GET("/users/register", middleware.GetLoginInfo, controllers.SignUpGET)
	router.POST("/users/register", controllers.SignUp)

	// Login routes
	router.GET("/users/login", middleware.GetLoginInfo, controllers.LoginGET)
	router.POST("/users/login", controllers.Login)

	// Users routes
	router.GET("/users/", middleware.Authenticate, controllers.UpdateUser)
	router.GET("/users/:userId", middleware.GetLoginInfo, controllers.UserProfileGET)
	// Saya menggunakan html untuk mengirim PUT dan DELETE jadi untuk router saya pake POST
	router.POST("/users/:userId", middleware.Authenticate, controllers.UpdateUser)

	// Photo routes
	router.GET("/photos/", middleware.GetLoginInfo, controllers.PhotosGET)
	router.POST("/photos/", middleware.Authenticate, controllers.PhotosPOST)
	// Saya menggunakan postman untuk handle PUT dan DELETE untuk photos
	router.PUT("/photos/:photoId", middleware.GetLoginInfo, controllers.PhotosUpdate)
	router.DELETE("/photos/:photoId", middleware.GetLoginInfo, controllers.PhotosDelete)

	router.Run()
}
