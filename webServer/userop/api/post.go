package api

import (
	"context"
	"net/http"
	"webServer/userop/forms"
	"webServer/userop/global"
	"webServer/userop/proto"

	"github.com/gin-gonic/gin"
)

type post struct{}

func NewPost() *post {
	return &post{}
}

func (p *post) List(c *gin.Context) {
	claims, err := getJwtToken(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": err.Error(),
		})
		return
	}

	var postListForm forms.PostListForm
	if err := c.ShouldBindQuery(&postListForm); err != nil {
		printValidateErrorTips(c, err)
		return
	}

	rsp, err := global.GrpcPostClient.QueryPostList(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.UserPostFilterRequest{
			UserId: int32(claims.ID),
			Page:   postListForm.Page,
			Limit:  postListForm.Limit,
		},
	)

	if err != nil {
		printGrpcError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": rsp,
	})
}

func (p *post) Add(c *gin.Context) {
	var postForm forms.PostForm
	if err := c.ShouldBindJSON(&postForm); err != nil {
		printValidateErrorTips(c, err)
		return
	}

	claims, err := getJwtToken(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": err.Error(),
		})
		return
	}
	_, err = global.GrpcPostClient.AddPost(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.UserPostRequest{
			UserId:  int32(claims.ID),
			Type:    postForm.Type,
			Subject: postForm.Subject,
			Message: postForm.Message,
			File:    postForm.File,
		},
	)
	if err != nil {
		printGrpcError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
