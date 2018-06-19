package handler

import (
	"github.com/labstack/echo"
	"github.com/backend/model"
	"github.com/globalsign/mgo/bson"
	"net/http"
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/smtp"
	"path/filepath"
)

type Account struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Host     string `json:"host"`
}

func ReadSecretJson() Account {
	// Read Secret JSON File
	absPath, _ := filepath.Abs("../src/github.com/backend/.secrets.json")
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

func SendActivationEmail(a Account, t string) {
	// 발신자의 SMTP 서버 Authentication
	auth := smtp.PlainAuth("", a.Email, a.Password, a.Host)

	// 발신자 주소: 썸띵모어 관리자
	from := a.Email

	// 수신자 주소: 가입자
	to := []string{t}

	// 본문
	ToHeader := "To: " + t + "\r\n"
	FromHeader := "From: " + a.Email + "\r\n"
	Subject := "Subject: 썸띵모어 회원 가입 인증 메일\r\n"
	Blank := "\r\n"
	body := "축하합니다!\r\n"
	msg := []byte(ToHeader + FromHeader + Subject + Blank + body)

	// 메일 전송
	err := smtp.SendMail(a.Host + ":587", auth, from, to, msg)

	// 에러 처리
	if err != nil {
		panic(err)
	}
}

func (h *Handler) SignUp(c echo.Context) (err error) {
	// Object bind
	u := &model.User{ID: bson.NewObjectId()}
	// Go 언어의 간단한 조건식:
	// 조건문 이전에 반드시 실행되는 구문을 세미콜론으로 구분해
	// if 문 안에서 실행하도록 한다
	if err = c.Bind(u); err != nil {
		return
	}

	// Validate
	if u.Email == "" || u.Password == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "이메일이나 패스워드가 입력되지 않았습니다",
		}
	}

	// Sending Email
	info := ReadSecretJson()
	go SendActivationEmail(info, u.Email) // go routine 을 사용한 비동기 처리

	// Save user
	db := h.DB.Clone()
	defer db.Close() // defer: 특정 문장이나 함수를 나중에 실행하게 해 줌
	if err = db.DB("st_more").C("users").Insert(u); err != nil {
		return
	}

	return c.JSON(http.StatusCreated, u)
}
