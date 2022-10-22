package api

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

var (
	accessKey = "T7Ucf6BG9omNCTwLwKmgn9PSLTyS1UZEPbL8B4qe"
	secretKey = "FeNKKsn4sA4K8jxMKzByk2UIIcj2NsnbtcX_uaeE"
	bucket    = "learn-gpmall"
)

// 自定义返回值结构体
type MyPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

// VerifyCallback 验证上传回调请求是否来自存储服务
// func VerifyCallback(mac *Mac, req *http.Request) (bool, error) {
// 	return mac.VerifyCallback(req)
// }

func main() {
	go func() {
		// server()
		web()
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	println("退出服务\n")
}

func web() {
	router := gin.Default()
	router.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			method := c.Request.Method
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "Content-type,Access-Token,X-CSRF-Token,Authorization,Token,x-token")
			c.Header("Access-Control-Allow-Methods", "POST,GET,PUT,OPTIONS,DELETE,PATCH")
			c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-type")
			c.Header("Access-Control-Allow-Credentials", "true")

			if method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
			}
		}
	}())
	uploadGroup := router.Group("/upload")
	{
		uploadGroup.GET("", func(c *gin.Context) {
			token := server()
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"token":   token,
				"message": "success",
			})
		})
		uploadGroup.POST("/callback", func(c *gin.Context) {
			mac := auth.New(accessKey, secretKey)
			if ok, err := mac.VerifyCallback(c.Request); !ok {
				panic(err.Error())
			}

			// myPutRet := make(map[string]interface{})
			myPutRet := MyPutRet{}
			c.BindJSON(&myPutRet)
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"data":    myPutRet,
				"message": "success",
			})
		})
	}
	router.Run("0.0.0.0:5355")
}

func server() string {

	// 简单上传凭证
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := auth.New(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	// fmt.Println(upToken)

	// // 设置上传凭证有效期
	// putPolicy = storage.PutPolicy{
	// 	Scope: bucket,
	// }
	// putPolicy.Expires = 7200 //示例2小时有效期

	// upToken = putPolicy.UploadToken(mac)
	// fmt.Println(upToken)

	// 覆盖上传凭证
	// 需要覆盖的文件名
	// keyToOverwrite := "qiniu.mp4"
	// putPolicy = storage.PutPolicy{
	// 	Scope: fmt.Sprintf("%s:%s", bucket, keyToOverwrite),
	// }
	// upToken = putPolicy.UploadToken(mac)
	// fmt.Println(upToken)

	// 自定义上传回复凭证
	// putPolicy = storage.PutPolicy{
	// 	Scope:      bucket,
	// 	ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
	// }
	// upToken = putPolicy.UploadToken(mac)
	// fmt.Println(upToken)

	// 带回调业务服务器的凭证(JSON方式)
	putPolicy = storage.PutPolicy{
		Scope:            bucket,
		CallbackURL:      "http://rryes5.natappfree.cc/upload/callback",
		CallbackBody:     `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
		CallbackBodyType: "application/json",
	}
	upToken = putPolicy.UploadToken(mac)
	// fmt.Println(upToken)

	// 带回调业务服务器的凭证（URL方式）
	// putPolicy = storage.PutPolicy{
	// 	Scope:        bucket,
	// 	CallbackURL:  "http://api.example.com/qiniu/upload/callback",
	// 	CallbackBody: "key=$(key)&hash=$(etag)&bucket=$(bucket)&fsize=$(fsize)&name=$(x:name)",
	// }
	// upToken = putPolicy.UploadToken(mac)
	// fmt.Println(upToken)

	// // 带数据处理的凭证
	// saveMp4Entry := base64.URLEncoding.EncodeToString([]byte(bucket + ":avthumb_test_target.mp4"))
	// saveJpgEntry := base64.URLEncoding.EncodeToString([]byte(bucket + ":vframe_test_target.jpg"))
	// //数据处理指令，支持多个指令
	// avthumbMp4Fop := "avthumb/mp4|saveas/" + saveMp4Entry
	// vframeJpgFop := "vframe/jpg/offset/1|saveas/" + saveJpgEntry
	// //连接多个操作指令
	// persistentOps := strings.Join([]string{avthumbMp4Fop, vframeJpgFop}, ";")
	// pipeline := "test"
	// putPolicy = storage.PutPolicy{
	// 	Scope:               bucket,
	// 	PersistentOps:       persistentOps,
	// 	PersistentPipeline:  pipeline,
	// 	PersistentNotifyURL: "http://api.example.com/qiniu/pfop/notify",
	// }
	// upToken = putPolicy.UploadToken(mac)
	// fmt.Println(upToken)
	return upToken
}
