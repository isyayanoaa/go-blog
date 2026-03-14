package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	cos "github.com/tencentyun/cos-go-sdk-v5"
	"github.com/isyayanoaa/go-blog/internal/config"
)

func cosClient() *cos.Client {
	cfg := config.C.COS
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", cfg.Bucket, cfg.Region))
	b := &cos.BaseURL{BucketURL: u}
	return cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	})
}

func Upload(c *gin.Context) {
	session := sessions.Default(c)
	if func() bool { v, ok := session.Get("is_admin").(int); return !ok || v != 1 }() {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file"})
		return
	}

	ext := filepath.Ext(file.Filename)
	key := fmt.Sprintf("uploads/%d%s", time.Now().UnixNano(), ext)

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "open file failed"})
		return
	}
	defer f.Close()

	client := cosClient()
	_, err = client.Object.Put(context.Background(), key, f, &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentLength: file.Size,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload to COS failed: " + err.Error()})
		return
	}

	cfg := config.C.COS
	publicURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", cfg.Bucket, cfg.Region, key)
	c.JSON(http.StatusOK, gin.H{"url": publicURL})
}
