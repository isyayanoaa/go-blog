package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RunPlayground(c *gin.Context) {
	session := sessions.Default(c)
	if func() bool { v, ok := session.Get("is_admin").(int); return !ok || v != 1 }() {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	var req struct {
		Code string `json:"code"`
	}
	c.BindJSON(&req)

	body, _ := json.Marshal(map[string]string{"body": req.Code})
	resp, err := http.Post("https://play.golang.org/compile", "application/json", bytes.NewReader(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	result, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", result)
}
