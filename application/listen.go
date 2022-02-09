package application

import (
	"es/controllers"

	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Listen() {
	posts := router.Group("posts")
	{
		posts.GET("", controllers.PostController.GetAllPosts)
		posts.POST("", controllers.PostController.AddPost)
		posts.GET(":id", controllers.PostController.GetPostById)
		posts.PATCH(":id", controllers.PostController.UpdatePost)
		posts.DELETE(":id", controllers.PostController.RemovePost)
	}
	router.Run(":3000")
}
