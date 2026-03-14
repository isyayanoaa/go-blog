package router

import (
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		// Auth
		auth := api.Group("/auth")
		{
			auth.GET("/github", nil)          // GitHub OAuth 跳转
			auth.GET("/github/callback", nil) // GitHub OAuth 回调
			auth.POST("/logout", nil)         // 登出
			auth.GET("/me", nil)              // 获取当前用户信息
		}

		// Posts
		posts := api.Group("/posts")
		{
			posts.GET("", nil)          // 文章列表
			posts.GET("/:id", nil)      // 文章详情
			posts.POST("", nil)         // 发布文章（管理员）
			posts.PUT("/:id", nil)      // 编辑文章（管理员）
			posts.DELETE("/:id", nil)   // 删除文章（管理员）
		}

		// Categories
		categories := api.Group("/categories")
		{
			categories.GET("", nil)         // 分类列表
			categories.GET("/:name", nil)   // 某分类下的文章
		}

		// Tags
		tags := api.Group("/tags")
		{
			tags.GET("", nil)        // 标签列表
			tags.GET("/:name", nil)  // 某标签下的文章
		}

		// Archives
		api.GET("/archives", nil) // 归档列表

		// Moments
		moments := api.Group("/moments")
		{
			moments.GET("", nil)         // 动态列表
			moments.GET("/:id", nil)     // 动态详情
			moments.POST("", nil)        // 发布动态（管理员）
			moments.DELETE("/:id", nil)  // 删除动态（管理员）
		}

		// Comments
		comments := api.Group("/comments")
		{
			comments.GET("", nil)         // 评论列表（by post_id 或 moment_id）
			comments.POST("", nil)        // 发布评论（登录用户）
			comments.DELETE("/:id", nil)  // 删除评论（管理员）
		}

		// Likes
		likes := api.Group("/likes")
		{
			likes.POST("", nil)    // 点赞
			likes.DELETE("", nil)  // 取消点赞
		}

		// Playground
		api.POST("/playground/run", nil) // 运行 Go 代码（管理员）

		// Upload
		api.POST("/upload", nil) // 上传文件（管理员）
	}

	return r
}
