package service

import (
	"github.com/isyayanoaa/go-blog/internal/model"
	"github.com/isyayanoaa/go-blog/internal/repo"
)

func UpsertGithubUser(githubID int64, username, avatar string, isAdmin bool) (uint, error) {
	return repo.UpsertUser(model.User{
		GithubID: githubID,
		Username: username,
		Avatar:   avatar,
		IsAdmin:  isAdmin,
	})
}
