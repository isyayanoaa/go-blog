package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/isyayanoaa/go-blog/internal/config"
	"github.com/isyayanoaa/go-blog/internal/service"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

func oauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.C.Github.ClientID,
		ClientSecret: config.C.Github.ClientSecret,
		RedirectURL:  config.C.Github.RedirectURL,
		Scopes:       []string{"user:email"},
		Endpoint:     githuboauth.Endpoint,
	}
}

func GithubLogin(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, oauthConfig().AuthCodeURL("state"))
}

func GithubCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
		return
	}
	token, err := oauthConfig().Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "oauth exchange failed"})
		return
	}
	resp, err := oauthConfig().Client(context.Background(), token).Get("https://api.github.com/user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var ghUser struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	}
	json.Unmarshal(body, &ghUser)

	isAdmin := ghUser.ID == config.C.Github.AdminGithubID
	userID, err := service.UpsertGithubUser(ghUser.ID, ghUser.Login, ghUser.AvatarURL, isAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	isAdminInt := 0
	if isAdmin {
		isAdminInt = 1
	}
	session := sessions.Default(c)
	session.Set("user_id", userID)
	session.Set("github_id", fmt.Sprintf("%d", ghUser.ID))
	session.Set("username", ghUser.Login)
	session.Set("avatar", ghUser.AvatarURL)
	session.Set("is_admin", isAdminInt)
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:5173")
}

func GetMe(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusOK, gin.H{"user": nil})
		return
	}
	isAdmin := false
	if v, ok := session.Get("is_admin").(int); ok && v == 1 {
		isAdmin = true
	}
	c.JSON(http.StatusOK, gin.H{"user": gin.H{
		"id":       userID,
		"username": session.Get("username"),
		"avatar":   session.Get("avatar"),
		"is_admin": isAdmin,
	}})
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
