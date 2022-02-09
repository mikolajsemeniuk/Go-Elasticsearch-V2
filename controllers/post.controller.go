package controllers

import (
	"es/inputs"
	"es/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var PostController IPostController = &postController{}

type IPostController interface {
	GetAllPosts(context *gin.Context)
	GetPostById(context *gin.Context)
	AddPost(context *gin.Context)
	UpdatePost(context *gin.Context)
	RemovePost(context *gin.Context)
}

type postController struct{}

func (*postController) GetAllPosts(context *gin.Context) {
	context.JSON(http.StatusOK, "GetAllPosts")
}

func (*postController) GetPostById(context *gin.Context) {
	param := context.Param("id")
	id, _ := uuid.Parse(param)

	payload, err := services.PostService.FindPostByID(id)
	if err != nil {
		// TODO: Add logger
		context.JSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": "error occured",
			"errors":  []string{err.Error()},
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"data":    payload,
		"message": "Post added",
		"errors":  []string{},
	})
}

func (*postController) AddPost(context *gin.Context) {
	input := inputs.Post{
		Title: "new title",
	}

	err := services.PostService.AddPost(input)
	if err != nil {
		// TODO: Add logger
		context.JSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": "error occured",
			"errors":  []string{err.Error()},
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"data":    nil,
		"message": "Post added",
		"errors":  []string{},
	})
}

func (*postController) UpdatePost(context *gin.Context) {
	param := context.Param("id")
	id, _ := uuid.Parse(param)
	input := inputs.Post{
		Title: "updated title",
	}

	err := services.PostService.UpdatePost(id, input)
	if err != nil {
		// TODO: Add logger
		context.JSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": "error occured",
			"errors":  []string{err.Error()},
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"data":    nil,
		"message": "Post updated",
		"errors":  []string{},
	})
}

func (*postController) RemovePost(context *gin.Context) {
	param := context.Param("id")
	id, _ := uuid.Parse(param)

	err := services.PostService.RemovePost(id)
	if err != nil {
		// TODO: Add logger
		context.JSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": "error occured",
			"errors":  []string{err.Error()},
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"data":    nil,
		"message": "Post removed",
		"errors":  []string{},
	})
}
