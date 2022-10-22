package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"

	"webServer/oss/global"
)

type qiniu struct {
	accessKey string
	secretKey string
	bucket    string
}

func NewQiniu() *qiniu {
	q := &qiniu{
		accessKey: global.ServerConfig.Kodo.AccessKey,
		secretKey: global.ServerConfig.Kodo.SecretKey,
		bucket:    global.ServerConfig.Kodo.Bucket,
	}
	return q
}

// 获取上传 token
func (q *qiniu) GetUpToken(c *gin.Context) {
	// 简单上传凭证
	putPolicy := storage.PutPolicy{
		Scope: q.bucket,
	}
	mac := auth.New(accessKey, secretKey)

	// 带回调业务服务器的凭证(JSON方式)
	putPolicy = storage.PutPolicy{
		Scope:            q.bucket,
		CallbackURL:      global.ServerConfig.Kodo.CallbackURL, //公网
		CallbackBody:     `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
		CallbackBodyType: "application/json",
	}
	upToken := putPolicy.UploadToken(mac)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"token":   upToken,
		"message": "success",
	})
}

// 自定义返回值结构体 ( CallbackBody )
type callbackBody struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

func (q *qiniu) Notify(c *gin.Context) {
	mac := auth.New(accessKey, secretKey)
	if ok, err := mac.VerifyCallback(c.Request); !ok {
		panic(err.Error())
	}

	// myPutRet := make(map[string]interface{})
	callbackBody := callbackBody{}
	c.BindJSON(&callbackBody)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    callbackBody,
		"message": "success",
	})
}

// 直传demo
func (q *qiniu) Kudo(c *gin.Context) {
	localFile := "/Users/jemy/Documents/github.png"
	bucket := "learn-gpmall"
	key := "github-x.png"

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	//自定义凭证有效期（示例2小时，Expires 单位为秒，为上传凭证的有效时间）
	putPolicy.Expires = 7200 //示例2小时有效期
	mac := qbox.NewMac(q.accessKey, q.secretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret.Key, ret.Hash)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
