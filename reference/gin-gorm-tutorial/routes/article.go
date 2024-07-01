package routes

import (
	"net/http"
	"strconv"
	"time"

	"../config"
	"../models"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func GetHome(c *gin.Context) {

	items := []models.Article{}
	config.DB.Find(&items)

	c.JSON( http.StatusOK, gin.H { 
		"status" : "berhasil ke halaman home",
		"data": items,
	})
}

func GetArticle(c *gin.Context) {
	slug := c.Param("slug")
	var item models.Article
	if err := config.DB.Where("slug = ?", slug).First(&item).Error; err != nil {
		c.JSON(404, gin.H {
			"status" : "error",
			"message" : "data tidak ditemukan",
		})
		c.Abort()
		return
	}

	c.JSON( http.StatusOK, gin.H { 
		"status" : "berhasil",
		"data": item,
	})
}

func PostArticle(c *gin.Context) {

	var oldItem models.Article
	var count int64
	slug := slug.Make(c.PostForm("title"))
	//if err := config.DB.Where("slug = ?", slug).First(&oldItem).Error; err != nil {
	if config.DB.Where("slug = ?", slug).First(&oldItem).Count(&count); count > 0 {
		slug = slug + "-" + strconv.FormatInt(time.Now().Unix(),10)
	}

	item := models.Article {
		Title : c.PostForm("title"),
		Desc : c.PostForm("desc"),
		Tag : c.PostForm("tag"),
		Slug : slug,
		UserID: uint(c.MustGet("jwt_user_id").(float64)),
	}

	config.DB.Create(&item)


	c.JSON( http.StatusOK, gin.H { 
		"status" : "berhasil ngepost",
		"data": item,
	})
}

func GetArticleByTag(c *gin.Context) {
	tag := c.Param("tag")
	items := []models.Article{}

	config.DB.Where("tag LIKE ?", "%" + tag + "%").Find(&items)
	c.JSON(200, gin.H{
		"data" : items,
	})

}


func UpdateArticle(c *gin.Context) {

	id := slug.Make(c.Param("id"))

	var item models.Article
	if err := config.DB.Where("id = ?", id).First(&item).Error; err != nil {
		c.JSON(404, gin.H {
			"status" : "error",
			"message" : "data tidak ditemukan",
		})
		c.Abort()
		return
	}

	if uint(c.MustGet("jwt_user_id").(float64)) != item.UserID{
		c.JSON(403, gin.H {
			"status" : "error",
			"message" : "data is forbidden",
		})
		c.Abort()
		return
	}

	config.DB.Model(&item).Where("id = ?" , id).Updates(models.Article{
			Title: c.PostForm("title"), 
			Desc: c.PostForm("desc"), 
			Tag: c.PostForm("tag"),
		})
		
	c.JSON( http.StatusOK, gin.H { 
		"status" : "berhasil",
		"data": item,
	})
}

func GetProfile(c *gin.Context) {
	var user models.User
	user_id := uint(c.MustGet("jwt_user_id").(float64))
	item := config.DB.Where("id = ? ", user_id).Preload("Articles", "user_id = ?", user_id).Find(&user)
	
	c.JSON( http.StatusOK, gin.H{ 
		"status" : "berhasil",
		"data": item,
	})
}

func DeleteArticle(c *gin.Context) {
	var article models.Article
	id := c.Param("id")

	config.DB.Where("id = ? ", id).Delete(&article)
	
	c.JSON( http.StatusOK, gin.H{ 
		"status" : "berhasil delete",
		"data": article,
	})
}
