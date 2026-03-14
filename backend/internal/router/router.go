package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	redisstore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/isyayanoaa/go-blog/internal/config"
	"github.com/isyayanoaa/go-blog/internal/handler"
	"github.com/isyayanoaa/go-blog/internal/middleware"
)

func Setup() *gin.Engine {
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Session (Redis)
	rc := config.C.Redis
	store, err := redisstore.NewStore(10,
		"tcp", fmt.Sprintf("%s:%d", rc.Host, rc.Port), "", rc.Password,
		[]byte("secret-key-32bytes-padded-123456"),
	)
	if err != nil {
		panic("redis store init failed: " + err.Error())
	}
	store.Options(sessions.Options{
		HttpOnly: true,
		// Secure: true, // 上线后开启（需要 HTTPS）
		MaxAge: 86400 * 7, // 7天
	})
	r.Use(sessions.Sessions("goblog_session", store))

	// Static uploads
	r.Static("/uploads", "./uploads")

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/github", handler.GithubLogin)
			auth.GET("/github/callback", handler.GithubCallback)
			auth.POST("/logout", handler.Logout)
			auth.GET("/me", handler.GetMe)
		}

		posts := api.Group("/posts")
		{
			// 公开：任何人可读
			posts.GET("", handler.GetPosts)
			posts.GET("/:id", handler.GetPost)
			// 写操作：仅 admin
			adminPosts := posts.Group("", middleware.AdminRequired())
			{
				adminPosts.POST("", handler.CreatePost)
				adminPosts.PUT("/:id", handler.UpdatePost)
				adminPosts.DELETE("/:id", handler.DeletePost)
			}
		}

		api.GET("/categories", handler.GetCategories)
		api.GET("/categories/:name", func(c *gin.Context) {
			c.Request.URL.RawQuery = "category=" + c.Param("name")
			handler.GetPosts(c)
		})
		api.GET("/tags", handler.GetTags)
		api.GET("/tags/:name", func(c *gin.Context) {
			c.Request.URL.RawQuery = "tag=" + c.Param("name")
			handler.GetPosts(c)
		})
		api.GET("/archives", handler.GetArchives)

		moments := api.Group("/moments")
		{
			// 公开：任何人可读
			moments.GET("", handler.GetMoments)
			// 写操作：仅 admin
			adminMoments := moments.Group("", middleware.AdminRequired())
			{
				adminMoments.POST("", handler.CreateMoment)
				adminMoments.DELETE("/:id", handler.DeleteMoment)
			}
		}

		comments := api.Group("/comments")
		{
			// 公开：任何人可读、可发评论（登录用户）
			comments.GET("", handler.GetComments)
			comments.POST("", middleware.LoginRequired(), handler.CreateComment)
			// 删评论：仅 admin
			comments.DELETE("/:id", middleware.AdminRequired(), handler.DeleteComment)
		}

		likes := api.Group("/likes")
		{
			likes.POST("", handler.AddLike)
			likes.DELETE("", handler.RemoveLike)
		}

		api.POST("/playground/run", handler.RunPlayground)
		api.POST("/upload", middleware.AdminRequired(), handler.Upload)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	return r
}
