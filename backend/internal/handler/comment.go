package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/isyayanoaa/go-blog/internal/dto"
	"github.com/isyayanoaa/go-blog/internal/service"
)

func GetComments(c *gin.Context) {
	comments, err := service.ListComments(c.Query("target_type"), c.Query("target_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

func CreateComment(c *gin.Context) {
	session := sessions.Default(c)
	userID, ok := session.Get("user_id").(int)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login required"}); return
	}
	var req dto.CreateCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	id, err := service.CreateComment(uint(userID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func DeleteComment(c *gin.Context) {
	if err := service.DeleteComment(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
