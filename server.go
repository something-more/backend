package main

import (
	// Default package
	"os"
	"time"
	"net/http"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	// Third Party package
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/globalsign/mgo"
	// User package
	"github.com/backend/handler"
)

// DBInfo 구조체 선언
type DBInfo struct {
	DBName   string `json:"db_name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// Echo instance
	e := echo.New()

	//-----------
	// Middleware
	//-----------

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//CORS WhiteList
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000", // master 에서는 변경할 것
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderAccept,
			echo.HeaderContentType,
			echo.HeaderXRequestedWith,
			echo.HeaderXCSRFToken,
			echo.HeaderAuthorization,
		},
		AllowCredentials: true,
	}))
	// XSRF Token
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookieSecure:   false, // master 에서는 변경할 것
		CookieHTTPOnly: false, // master 에서는 변경할 것
	}))
	// JWT
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(handler.Key), // "secret"
		Skipper: func(c echo.Context) bool {
			// 인증 메서드의 경우 authentication 을 건너뛴다
			if c.Path() == "/" ||
				c.Path() == "/assets/*" ||
				c.Path() == "/admin/" ||
				c.Path() == "/sign-up/" ||
				c.Path() == "/sign-in/" ||
				c.Path() == "/activate/" ||
				c.Path() == "/authors/" ||
				c.Path() == "/authors/:author_id" ||
				c.Path() == "/authors/count/:author_id" ||
				c.Path() == "/story/client/" ||
				c.Path() == "/story/view/:story_id" ||
				c.Path() == "/board/list/" ||
				c.Path() == "/board/count/" ||
				c.Path() == "/board/view/:board_id" ||
				c.Path() == "/notice/list/" ||
				c.Path() == "/notice/count/" ||
				c.Path() == "/notice/view/:notice_id" {
				return true
			}
			return false
		},
	}))

	//-----------
	// Databases
	//-----------

	// Read Secreat JSON File
	absPath, _ := filepath.Abs("../src/github.com/backend/.secrets_db.json")
	jsonFile, err := os.Open(absPath)
	defer jsonFile.Close()
	if err != nil {
		e.Logger.Fatal(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Bind to DBInfo struct
	var data DBInfo
	json.Unmarshal(byteValue, &data)

	// Database connection
	info := &mgo.DialInfo{
		Addrs:    []string{"localhost"},
		Timeout:  60 * time.Second,
		Database: data.DBName,
		Username: data.Username,
		Password: data.Password,
	}

	db, err := mgo.DialWithInfo(info)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Create indices
	// 인덱스 값으로 email 을 사용하며, 그 값은 고유하다
	if err = db.Copy().DB(handler.DBName).C(handler.USER).EnsureIndex(mgo.Index{
		Key:    []string{"$text:email", "$text:nickname"},
		Unique: true,
	}); err != nil {
		log.Fatal(err)
	}

	//---------------
	// Route & Server
	//---------------

	// Initialize handler
	h := &handler.Handler{DB: db}

	// Route: Static
	e.Static("/assets", "assets") // 정적 파일

	// Route: Index
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "섬띵모어 API 서버\n")
	})

	// Route: User
	e.POST("/sign-up/", h.SignUpNormal)    // 회원 가입
	e.POST("/admin/", h.SignUpAdmin)       // 관리자 회원 가입
	e.GET("/activate/", h.Activate)        // 이메일 회원 활성화
	e.POST("/sign-in/", h.SignIn)          // 로그인
	e.PATCH("/patch/", h.PatchPassword)    // 비밀번호 수정
	e.PATCH("/nickname/", h.PatchNickname) // 닉네임 수정
	e.DELETE("/destroy/", h.DestroyUser)   // 회원 탈퇴

	// Route: Admin
	e.GET("/users/", h.ListUsers)                   // 전체 유저 리스트
	e.PATCH("/users/:user_email", h.UpdateUserAuth) // 유저

	// Route: Author
	e.GET("/authors/", h.ListAuthors)                      // 필진 리스트
	e.GET("/authors/:author_id", h.ListStoryAuthor)        // 필진 스토리 리스트
	e.GET("/authors/count/:author_id", h.CountStoryAuthor) // 필진 스토리 갯수

	// Route: Story
	e.POST("/story/", h.CreateStory)                          // 스토리 생성
	e.GET("/story/", h.ListStory)                             // 스토리 리스트
	e.GET("/story/client/", h.ClientListStory)                // 클라이언트 스토리 리스트
	e.GET("/story/count/", h.CountStory)                      // 스토리 총 갯수
	e.GET("/story/view/:story_id", h.RetrieveStory)           // 스토리 디테일
	e.PATCH("/story/:story_id", h.PatchStory)                 // 스토리 수정
	e.PATCH("/story/publish/:story_id", h.ChangePublishStory) // 스토리 발행 상태 변경
	e.DELETE("/story/:story_id", h.DestroyStory)              // 스토리 삭제

	// Route: Board
	e.POST("/board/", h.CreateBoard)                // 자유게시판 글 생성
	e.GET("/board/list/", h.ListBoard)              // 자유게시판 글 목록
	e.GET("/board/count/", h.CountBoard)            // 자유게시판 글 갯수
	e.GET("/board/view/:board_id", h.RetrieveBoard) // 자유게시판 글 보기
	e.PATCH("/board/:board_id", h.PatchBoard)       // 자유게시판 글 수정
	e.DELETE("/board/:board_id", h.DestroyBoard)    // 자유게시판 글 삭제

	// Route: Notice
	e.POST("/notice/", h.CreateNotice)                 // 공지사항 글 생성
	e.GET("/notice/list/", h.ListNotice)               // 공지사항 글 목록
	e.GET("/notice/count/", h.CountNotice)             // 공지사항 글 갯수
	e.GET("/notice/view/:notice_id", h.RetrieveNotice) // 공지사항 글 보기
	e.PATCH("/notice/:notice_id", h.PatchNotice)       // 공지사항 글 수정
	e.DELETE("/notice/:notice_id", h.DestroyNotice)    // 공지사항 글 삭제

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
