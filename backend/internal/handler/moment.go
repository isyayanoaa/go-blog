package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isyayanoaa/go-blog/internal/dto"
	"github.com/isyayanoaa/go-blog/internal/service"
)

func GetMoments(c *gin.Context) {
	moments, err := service.ListMoments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"moments": moments})
}

func CreateMoment(c *gin.Context) {
	var req dto.CreateMomentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	id, err := service.CreateMoment(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func DeleteMoment(c *gin.Context) {
	if err := service.DeleteMoment(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
