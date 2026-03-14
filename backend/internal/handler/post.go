package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isyayanoaa/go-blog/internal/dto"
	"github.com/isyayanoaa/go-blog/internal/service"
)

func GetPosts(c *gin.Context) {
	posts, err := service.ListPosts(c.Query("category"), c.Query("tag"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func GetPost(c *gin.Context) {
	post, err := service.GetPost(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"}); return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}

func CreatePost(c *gin.Context) {
	var req dto.CreatePostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	id, err := service.CreatePost(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func UpdatePost(c *gin.Context) {
	var req dto.UpdatePostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	if err := service.UpdatePost(c.Param("id"), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func DeletePost(c *gin.Context) {
	if err := service.DeletePost(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func GetCategories(c *gin.Context) {
	cats, _ := service.ListCategories()
	c.JSON(http.StatusOK, gin.H{"categories": cats})
}

func GetTags(c *gin.Context) {
	tags, _ := service.ListTags()
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

func GetArchives(c *gin.Context) {
	arcs, _ := service.ListArchives()
	c.JSON(http.StatusOK, gin.H{"archives": arcs})
}
