package controllers

import (
	"net/http"
	"strconv"
	"task/api/database"
	"task/api/models"
	"time"

	"github.com/gin-gonic/gin"
)

func PhotosGET(c *gin.Context) {
	user, _ := c.Get("user")
	empty := models.User{}
	if user == empty {
		c.HTML(http.StatusOK, "photos.html", nil)
		return
	}
	c.HTML(http.StatusOK, "photos.html", user)

}

func PhotosPOST(c *gin.Context) {
	if c.Request.Method == "POST" {
		var body struct {
			Title    string `json:"title" form:"title"`
			Captions string `json:"captions" form:"captions"`
		}
		// Bind body
		if c.ShouldBind(&body) != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "Failed to read body",
			})
			return
		}
		// Get user information
		user, _ := c.Get("user")
		empty := models.User{}
		if user == empty {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "User not found"})
			return
		}

		var latestID int
		database.DB.Model(&models.Picture{}).Table("pictures").Select("MAX(id) AS latest_id").Scan(&latestID)
		latestID = latestID + 1
		newID := "/photos/" + strconv.Itoa(latestID)
		newPhoto := models.Picture{
			Title:     body.Title,
			Caption:   body.Captions,
			PhotoUrl:  newID,
			UserID:    user.(models.User).ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		result := database.DB.Create(&newPhoto)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": result.Error})
			return
		}
		msg := "Photo created successfully"
		c.JSON(http.StatusOK, gin.H{"message": msg})
	} else {
		c.HTML(http.StatusOK, "photos.html", nil)
	}

}

func PhotosUpdate(c *gin.Context) {
	if c.Request.Method == "PUT" {
		user, _ := c.Get("user")
		empty := models.User{}
		if user == empty {
			c.HTML(http.StatusOK, "photos.html", nil)
			return
		}
		var body struct {
			Title   string `json:"title" form:"title"`
			Caption string `json:"caption" form:"caption"`
		}

		if c.ShouldBind(&body) != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "Failed to read body",
			})
			return
		}
		// Get photo id from url
		photoId := c.Param("photoId")
		var photo models.Picture
		database.DB.First(&models.Picture{}, "id = ?", photoId).Scan(&photo)
		userID := user.(models.User).ID
		if photo.UserID != userID {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "You are not authorized to update this photo"})
			return
		}
		database.DB.Model(&photo).Where("id = ?", photoId).Select("title", "caption", "updated_at").Updates(models.Picture{Title: body.Title, Caption: body.Caption, UpdatedAt: time.Now()})

		msg := "Photo" + photo.Title + " updated successfully"
		c.JSON(http.StatusOK, gin.H{"message": msg})

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request"})

	}
}

func PhotosDelete(c *gin.Context) {
	if c.Request.Method == "DELETE" {
		// Get user information
		user, _ := c.Get("user")
		userID := user.(models.User).ID
		// Get photoId from url
		photoId := c.Param("photoId")
		var photo models.Picture
		database.DB.First(&models.Picture{}, "id = ?", photoId).Scan(&photo)
		//Compare user id with photo user id
		if photo.UserID != userID {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "You are not authorized to update this photo"})
			return
		}
		database.DB.Unscoped().Model(&photo).Where("id =?", photoId).Delete(&models.Picture{})
		c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request"})
	}
}
