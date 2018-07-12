package utility

import (
	// Default package
	"os"
	"fmt"
	"bytes"
	"net/smtp"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"html/template"
	// Third Party package
	"github.com/labstack/echo"
	// User package
	"github.com/backend/model"
)

type Account struct {
	Email    string // 발신자
	Password string // 패스워드
	Host     string // SMTP 서버
}

type Request struct {
	From    string   // 발신자
	To      []string // 수신자
	Subject string   // 제목
	Body    string   // 내용
}

type TemplateData struct {
	Title    string // 이메일 주소
	Nickname string // 유저 닉네임
	URL      string // Activate 주소
}

func ReadSecretJson() Account {
	// Read Secret JSON File
	absPath, _ := filepath.Abs("./secrets/.secrets_email.json")
	jsonFile, err := os.Open(absPath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	// Account 구조체 변수 선언
	var account Account

	// JSON 파일을 읽어 byte 값을 account 변수의 주소에 넣는다
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &account)

	return account
}

func (r *Request) ParseTemplate(name string, data interface{}) (err error) {
	t, err := template.ParseFiles(name)
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return
	}
	r.Body = buf.String()
	return
}

func (r *Request) SendEmail(host string, auth smtp.Auth) (bool, error) {
	// MIME Type 정의
	MIME := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	// To 헤더
	ToHeader := "To: " + r.To[0] + "\r\n"
	// From 헤더
	FromHeader := "From: " + r.From + "\r\n"
	// 제목
	SubjectHeader := "Subject: " + r.Subject + "\r\n"
	// 전체 메시지
	msg := []byte(ToHeader + FromHeader + SubjectHeader + MIME + "\r\n" + r.Body)
	// 발신자 SMTP 서버
	addr := host + ":587"

	// 이메일 전송
	if err := smtp.SendMail(addr, auth, r.From, r.To, msg); err != nil {
		return false, err
	}
	return true, nil
}

func SendActivationEmail(c echo.Context, u *model.User) (err error) {
	// Secret json 읽기
	s := ReadSecretJson()

	// 발신자의 SMTP 서버 Authentication
	auth := smtp.PlainAuth("", s.Email, s.Password, s.Host)

	// Request 객체 생성
	r := &Request{
		From:    s.Email,            // 발신자 주소: 썸띵모어 관리자
		To:      []string{u.Email},  // 수신자 주소: 가입자
		Subject: "썸띵모어 회원 가입 인증 메일", // 메일 제목
	}

	// TemplateData 객체 생성
	d := &TemplateData{
		Title:    r.Subject,
		Nickname: u.Nickname,
		URL:      c.Scheme() + "://" + c.Request().Host + "/activate/" + u.Email,
	}

	// Template File 경로 생성
	templatePath, _ := filepath.Abs("./templates/activate_account.html")

	if err = r.ParseTemplate(templatePath, d); err == nil {
		r.SendEmail(s.Host, auth)
	}
	return
}
