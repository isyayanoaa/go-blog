package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/isyayanoaa/go-blog/internal/dto"
	"github.com/isyayanoaa/go-blog/internal/service"
)

func AddLike(c *gin.Context) {
	session := sessions.Default(c)
	userID, ok := session.Get("user_id").(int)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login required"}); return
	}
	var req dto.LikeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	if err := service.AddLike(uint(userID), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func RemoveLike(c *gin.Context) {
	session := sessions.Default(c)
	userID, ok := session.Get("user_id").(int)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login required"}); return
	}
	var req dto.LikeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	if err := service.RemoveLike(uint(userID), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
