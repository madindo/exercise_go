package main

import (
	"./config"
	"./middleware"
	"./routes"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

func main() {
	config.InitDB()
	gotenv.Load()

	r := gin.Default()

	v1 := r.Group("/api/v1/") 
	{
		v1.GET("auth/:provider", routes.RedirectHandler)
		v1.GET("auth/:provider/callback", routes.CallbackHandler)
		v1.GET("check", middleware.IsAuth(), routes.CheckToken)
		v1.GET("profile", middleware.IsAuth(), routes.GetProfile)

		v1.GET("/article/:slug", routes.GetArticle)
		articles := v1.Group("/articles")
		{
			articles.GET("/", routes.GetHome)
			articles.POST("/", middleware.IsAuth(),  routes.PostArticle)
			articles.GET("/tag/:tag", routes.GetArticleByTag)
			articles.PUT("/update/:id", middleware.IsAuth(),  routes.UpdateArticle)
			articles.DELETE("/delete/:id", middleware.IsAdmin(),  routes.DeleteArticle)
		}
		
	}
	
	r.Run() 
}