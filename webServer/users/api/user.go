package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"webServer/users/forms"
	"webServer/users/global"
	"webServer/users/global/response"
	"webServer/users/middlewares"
	"webServer/users/models"
	"webServer/users/proto"
)

type user struct{}

func NewUser() *user {
	t := &user{}
	return t
}

// redis connect
func connect_redis() (*redis.Client, context.Context) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var c = context.Background()

	// 心跳
	pong, err := rdb.Ping(c).Result()
	if err != nil {
		zap.S().Debug(pong, err)
	}

	return rdb, c
}

/*
	创建用户
*/
func (u user) CreateUser(c *gin.Context) {
	var registerForm forms.RegisterForm
	if err := c.ShouldBind(&registerForm); err != nil {
		HandleValidateHttp(err, c)
		return
	}

	// validate captcha
	if !store.Verify(registerForm.CaptchaId, registerForm.VerifyValue, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   2,
			"errmsg": "验证码有误",
		})
		return
	}

	// 此处需单独开接口发送(如 captcha)
	go func() {
		if err := NewEmail().SendEmailTo("mclgo2020@126.com", "用户注册验证码", 6); err != nil {
			//验证码发送失败
			println("Email 验证码发送失败")
		}
	}()

	// 验证邮件验证码
	// rdb, ctx := connect_redis()
	// code, err := rdb.Get(ctx, registerForm.Mobile).Result()
	// if err == nil {
	// 	println(code)
	// 	if registerForm.code != code {
	// 		panic("验证失败")
	// 	}
	// }

	rsp, err := global.UserClient.CreateUser(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.CreateUserRequest{
			Mobile:     registerForm.Mobile,
			NickName:   registerForm.Nick_Name,
			Password:   registerForm.Password,
			Repassword: registerForm.Repassword,
			Birthday:   registerForm.Birthday,
			Gender:     registerForm.Gender,
			Role:       registerForm.Role,
			Avatar:     registerForm.Avatar,
			Desc:       registerForm.Desc,
			Country:    registerForm.Country,
			Provice:    registerForm.Provice,
			City:       registerForm.City,
			Area:       registerForm.Area,
			Address:    registerForm.Address,
		})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	token, err := middlewares.NewJWT().GenerateToken(models.CustomClaims{
		ID:          uint(rsp.Id),
		NickName:    rsp.NickName,
		AuthorityId: uint(rsp.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               //签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //过期时间
			Issuer:    "imooc",
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   401,
			"errmsg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    rsp,
		"x-token": token,
		"message": "success",
	})
}

/*
	登陆
	jwt
*/
func (u user) Login(c *gin.Context) {
	var checkLoginForm forms.CheckLoginForm
	if err := c.ShouldBind(&checkLoginForm); err != nil {
		HandleValidateHttp(err, c)
		return
	}

	// 1. get user info
	rsp, err := global.UserClient.GetUserInfo(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.MobileRequest{
			Mobile: checkLoginForm.Username,
		},
	)
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	// 2. check password
	if rsp.Password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": "无效用户：密码未设置",
		})
		return
	}
	rsp1, err := global.UserClient.CheckPassword(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.CheckPasswordRequest{
			Password:          checkLoginForm.Password,
			EncryptedPassword: rsp.Password,
		},
	)
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	if rsp1.Success == false {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": "密码不正确",
		})
		return
	}

	// 3. jwt鉴权
	token, err := middlewares.NewJWT().GenerateToken(models.CustomClaims{
		ID:          uint(rsp.Id),
		NickName:    rsp.NickName,
		AuthorityId: uint(rsp.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               //签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //过期时间
			Issuer:    "gpmall",
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   401,
			"errmsg": err.Error(),
		})
		return
	}

	// success response
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"status":  rsp1.Success,
		"x-token": token,
		"message": "login success",
	})
}

/*
	获取用户信息 by mobile
*/
func (u user) GetUserInfo(c *gin.Context) {
	var mobileForm forms.MobileForm
	if err := c.ShouldBind(&mobileForm); err != nil {
		HandleValidateHttp(err, c)
		return
	}

	rsp, err := global.UserClient.GetUserInfo(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.MobileRequest{
			Mobile: mobileForm.Mobile,
		},
	)
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": rsp,
		// "message": "success",
	})
}

/*
	获取用户信息 by id
*/
func (u user) GetUserById(c *gin.Context) {
	// 1. 验证参数
	var uidForm forms.UidForm
	if err := c.ShouldBindUri(&uidForm); err != nil {
		HandleValidateHttp(err, c)
		return
	}

	rsp, err := global.UserClient.GetUserById(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.UidRequest{
			Uid: uidForm.Uid,
		},
	)
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": rsp,
		// "message": "success",
	})
}

/*
	获取用户列表
*/
func (u user) GetUserList(c *gin.Context) {
	var userListForm forms.UserListForm
	if err := c.ShouldBind(&userListForm); err != nil {
		HandleValidateHttp(err, c)
		return
	}

	rsp, err := global.UserClient.GetUserList(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.PageRequest{
			Page:  userListForm.Page,
			Limit: userListForm.Limit,
		},
	)
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	//2. provide api server for web browser
	// c.ShouldBind()
	data := make([]interface{}, 0)
	for _, v := range rsp.Data {
		row := response.UserResponse{
			Id:       v.Id,
			NickName: v.NickName,
			Birthday: v.Birthday,
			Gender:   v.Gender,
			Mobile:   v.Mobile,
			Password: v.Password,
			Role:     v.Role,
			Avatar:   v.Avatar,
			Desc:     v.Desc,
			Country:  v.Country,
			Provice:  v.Provice,
			City:     v.City,
			Area:     v.Area,
			Address:  v.Address,
		}
		data = append(data, row)
	}

	// get jwt token
	claims, _ := c.Get("claims")
	claim, _ := claims.(*models.CustomClaims)

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"data":    data,
		"currid":  claim.ID,
		"message": "success",
	})
}

/*
更新用户信息
*/
func (u user) UpdateUserInfo(c *gin.Context) {
	var updateUserRequest forms.UpdateUserForm
	if err := c.ShouldBind(&updateUserRequest); err != nil {
		HandleValidateHttp(err, c)
		return
	}

	_, err := global.UserClient.UpdateUserInfo(
		context.WithValue(context.Background(), "ginContext", c),
		&proto.UpdateUserRequest{
			Uid: updateUserRequest.Uid,
			Data: &proto.CreateUserRequest{
				Mobile:   updateUserRequest.Mobile,
				NickName: updateUserRequest.Nick_Name,
				Birthday: updateUserRequest.Birthday,
				Gender:   updateUserRequest.Gender,
				Role:     updateUserRequest.Role,
				Avatar:   updateUserRequest.Avatar,
				Desc:     updateUserRequest.Desc,
				Country:  updateUserRequest.Country,
				Provice:  updateUserRequest.Provice,
				City:     updateUserRequest.City,
				Area:     updateUserRequest.Area,
				Address:  updateUserRequest.Address,
			},
		})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	//2.输出
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// 格式化错误信息
func ErrToJson(errs validator.ValidationErrors) map[string]string {
	tips := errs.Translate(global.Trans)
	data := make(map[string]string)
	var t []string
	for k, v := range tips {
		t = strings.Split(k, ".")
		data[t[len(t)-1]] = v
	}

	return data
}

// 验证字段
func HandleValidateHttp(err error, c *gin.Context) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": ErrToJson(errs),
		})
	} else {
		zap.S().Info("数据类型错误[err.(validator.ValidationErrors)]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"errmsg": "数据类型有误",
		})
	}
}

// 将grpc的code转换成http的状态码
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	// 将grpc的code转换成http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"code": 1, "errmsg": e.Message()})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "errmsg": e.Message()})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{"code": 1, "errmsg": e.Message()})
			case codes.Unauthenticated:
				c.JSON(http.StatusUnauthorized, gin.H{"code": 1, "errmsg": e.Message()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "errmsg": e.Message()})
			}
			zap.S().Debugf("grpc错误: %s => %s", e.Code(), e.Message())
		}
	}
}
