package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"webServer/order/models"
	"webServer/order/utils/translator"
)

func getJwtInfo(c *gin.Context) (*models.MyCustomClaims, error) {
	claims, ok := c.Get("claims")
	if !ok {
		return nil, errors.New("")
	}
	if session, ok := claims.(*models.MyCustomClaims); ok {
		return session, nil
	}
	return nil, errors.New("")
}

func printGrpcError(c *gin.Context, err error) {
	if s, ok := status.FromError(err); ok {
		switch s.Code() {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "errmsg": s.Message()})
		case codes.Unauthenticated:
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "errmsg": s.Message()})
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "errmsg": s.Message()})
		case codes.PermissionDenied:
			c.JSON(http.StatusNotAcceptable, gin.H{"code": 406, "errmsg": s.Message()})
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "errmsg": s.Message()})
		case codes.Unavailable:
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "errmsg": s.Message()})
		default:
			c.JSON(http.StatusBadGateway, gin.H{"code": 502, "errmsg": s.Message()})
		}
	} else {
		c.JSON(http.StatusBadGateway, gin.H{"code": 502, "errmsg": s.Message()})
	}
}

func printValidateError(c *gin.Context, err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": translator.FormatValidationError(errs),
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   2,
			"errmsg": "请求参数类型错误",
		})
	}
}
