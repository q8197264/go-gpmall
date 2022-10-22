package api

import (
	"fmt"
	"log"
	"math/rand"
	"mime"
	"strings"
	"sync"
	"time"

	"gopkg.in/gomail.v2"
)

type email struct {
	username string
	authCode string
	host     string
	port     int

	wg sync.WaitGroup
}

func NewEmail() *email {
	//定义邮箱服务器连接信息，如果是网易邮箱 pass填密码，qq邮箱填授权码
	return &email{
		username: "mclgo2021@126.com",
		authCode: "RUHCSQDKJZCJMNHV",
		host:     "smtp.126.com",
		port:     465,
		wg:       sync.WaitGroup{},
	}
}

/**
	发送 email
	@pargam mailTo string 收件人 如："mclgo2020@126.com",
	@param subject string 邮件主题
	@param length int 验证码长度
**/
func (e *email) SendEmailTo(mailTo string, subject string, length int) error {
	rand.Seed(time.Now().Unix())

	// sms := NewSms()

	// 定义收件人
	// mailTo := []string{
	// 	"mclgo2020@126.com",
	// 	// "*****@qq.com",
	// }

	// 邮件主题
	// subject := "Hello,Go Mail"

	// length := 6

	e.wg.Add(1)
	err := e.sendMail(mailTo, subject, length)
	e.wg.Wait()

	return err
}

/**
	批量发送 email
	@pargam mailTo []string 收件人(可多个) 如：[]string{"mclgo2020@126.com","xxx"}
	@param subject string 邮件主题
	@param length int 验证码长度
**/
func (e *email) SendBatchEmailTo(mailTo []string, subject string, length int) {
	rand.Seed(time.Now().Unix())

	e.wg.Add(len(mailTo))
	for _, mail := range mailTo {
		go e.sendMail(mail, subject, length)
	}
	e.wg.Wait()
}

// 三方包：发送邮件
func (e *email) sendMail(mailTo string, subject string, length int) error {
	defer e.wg.Done()

	g := gomail.NewMessage()
	g.SetHeader("From", mime.QEncoding.Encode("UTF-8", "闭关自学上岸中心")+"<"+e.username+">")
	g.SetHeader("To", mailTo)
	g.SetHeader("Subject", subject)
	g.SetBody("text/html", e.generatebody(length))

	/*
		name := "附件.txt"
		m.Attach("/tmp/foo.txt",
			gomail.Rename(name),
		)
	*/

	d := gomail.NewDialer(e.host, e.port, e.username, e.authCode)
	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}//这里的用户名和密码指的是能够登录该邮箱的邮箱地址和密码如果端口号为465的话，自动开启SSL，这个时候需要指定TLSConfig
	err := d.DialAndSend(g)
	if err != nil {
		log.Fatalln(" SendTo:", mailTo, "mclgo2020@126.com", "Send Email Failed!Err:", err)
		return err
	} else {
		log.Println("Send To:", mailTo, "mclgo2020@126.com", "Send Email Successfully!")
	}
	return nil
}

// 邮件内容
func (e *email) generatebody(length int) string {
	// 邮件主题
	html := `<h1>Hello From Go Mail 验证码是 %s </h1>
		<p><a href="http://www.runoob.com">重设密码</a></p>
	`
	code := e.generateSmsCode(length)
	html = fmt.Sprintf(html, code)
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "", // no password set
	// 	DB:       0,  // use default DB
	// })
	// var c context.Context
	// if err := rdb.Set(c, "", code, 2).Err(); err != nil {
	// 	panic(err)
	// }

	return html
}

func (e *email) generateSmsCode(length int) string {
	// b := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	// r := len(b)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < length; i++ {
		// fmt.Fprintf(&sb, "%d", b[rand.Intn(r)])
		fmt.Fprintf(&sb, "%d", rand.Intn(10))
	}
	return sb.String()
}
