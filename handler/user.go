package handler

import (
	// Default package
	"os"
	"fmt"
	"time"
	"net/http"
	"net/smtp"
	"path/filepath"
	"crypto/sha256"
	"io/ioutil"
	"encoding/json"
	"encoding/hex"
	// Third-party package
	"github.com/labstack/echo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/dgrijalva/jwt-go"
	// User package
	"github.com/backend/model"
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

func HashPassword(p string) string {
	// 패스워드를 SHA-256 알고리즘으로 암호화
	// string 타입인 rawPassword 를 byte 배열에 삽입
	rawPassword := []byte(p)
	// SHA-256 알고리즘으로 Hash
	sum := sha256.Sum256(rawPassword)
	// sum 배열 요소 전체를 호출해 string 타입으로 변경
	return hex.EncodeToString(sum[:])
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

	// Hash Password
	newPassword := HashPassword(u.Password)

	// hash 된 패스워드를 원래 자리에 대입
	u.Password = newPassword

	// Save user
	db := h.DB.Clone()
	defer db.Close() // defer: 특정 문장이나 함수를 나중에 실행하게 해 줌
	if err = db.DB("st_more").C("users").Insert(u); err != nil {
		return
	}

	// Sending Email
	info := ReadSecretJson()
	go SendActivationEmail(info, u.Email) // go routine 을 사용한 비동기 처리

	return c.JSON(http.StatusCreated, u)
}

func (h *Handler) Activate(c echo.Context) (err error) {
	// Object bind
	// Signup 과 달리 비어 있는 객체를 생성
	u := new(model.User)
	if err = c.Bind(u); err != nil {
		return
	}

	// Find user
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("st_more").C("users").
		Find(bson.M{"email": u.Email}).One(u); err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "이메일이 올바르지 않습니다",
			}
			return
		}
	}

	// Active user
	u.IsActive = true

	// 메인 페이지로 리다이렉트
	return c.Redirect(http.StatusMovedPermanently, "http://localhost:3000")
}

func (h *Handler) Login(c echo.Context) (err error) {
	// Object bind
	// 비어 있는 객체 생성
	u := new(model.User)
	if err = c.Bind(u); err != nil {
		return
	}

	// Hash password
	// 로그인 시 입력한 패스워드를 해쉬해서 DB 안에 있는 패스워드와 비교한다
	comparePassword := HashPassword(u.Password)

	// Find user
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("st_more").C("users").
		Find(bson.M{"email": u.Email, "password": comparePassword}).One(u); err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{Code: http.StatusUnauthorized,
				Message: "이메일이나 패스워드가 올바르지 않습니다"}
		}
		return
	}

	//-----
	// JWT
	//-----

	// Create token
	// HS256 알고리즘으로 인코딩
	// 단방향 암호화 알고리즘인 RS256과 달리 양방향 암호화 알고리즘이므로 디코딩이 가능함
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	// 유저 정보를 담는다
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = u.ID
	claims["email"] = u.Email
	claims["isActive"] = u.IsActive
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // 토큰 유효시간: 72시간

	// 토큰 인코딩 및 response 에 추가하기
	// signing key 로 핸들러에 정의해 둔 Key 상수를 사용
	u.Token, err = token.SignedString([]byte(Key))
	if err != nil {
		return err
	}

	// 최종적으로는 암호화된 토큰만 전송한다
	return c.JSON(http.StatusOK, u.Token)
}