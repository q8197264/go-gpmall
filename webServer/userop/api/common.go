package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"webServer/userop/models"
	Trans "webServer/userop/utils/translator"
)

func getJwtToken(c *gin.Context) (customClaims *models.MyCustomClaims, e error) {
	claims, ok := c.Get("claims")
	if !ok {
		return customClaims, errors.New("fail get claims ")
	}

	if customClaims, ok = claims.(*models.MyCustomClaims); ok {
		return customClaims, nil
	}

	return customClaims, errors.New("fail models.MyCustomClaims")
}

func printValidateErrorTips(c *gin.Context, errs error) {
	if err, ok := errs.(validator.ValidationErrors); ok {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": Trans.TranslateValidateErrors(err),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": "类型错误",
		})
	}
}

func printGrpcError(c *gin.Context, err error) {
	zap.S().Warnf("grpc错误:%s", err.Error())
	if s, ok := status.FromError(err); ok {
		switch s.Code() {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "errmsg": fmt.Sprintf("记录不存在:%s", s.Message())})
		case codes.Unauthenticated:
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "errmsg": fmt.Sprintf("未授权:%s", s.Message())})
		case codes.PermissionDenied:
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "errmsg": fmt.Sprintf(":%s", s.Message())})
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "errmsg": s.Message()})
		case codes.Unavailable:
			c.JSON(http.StatusServiceUnavailable, gin.H{"code": 503, "errmsg": fmt.Sprintf("连接失败:%s", s.Message())})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"code": s.Code(), "errmsg": fmt.Sprintf("内部错误:%s", s.Message())})
		}
	} else {
		c.JSON(http.StatusBadGateway, gin.H{"code": 502, "errmsg": s.Message()})
	}
}
