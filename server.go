package main

import (
	// Default package
	"net/http"
	// Third Party package
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/globalsign/mgo"
	// User package
	"github.com/backend/handler"
)

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
		CookieSecure:	false, // master 에서는 변경할 것
		CookieHTTPOnly: false, // master 에서는 변경할 것
	}))
	// JWT
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(handler.Key), // "secret"
		Skipper: func(c echo.Context) bool {
			// 인증 메서드의 경우 authentication 을 건너뛴다
			if c.Path() == "/" ||
				c.Path() == "/admin/" ||
				c.Path() == "/sign-up/" ||
				c.Path() == "/sign-in/" ||
				c.Path() == "/activate/" {
				return true
			}
			return false
		},
	}))

	//-----------
	// Databases
	//-----------

	// Database connection
	db, err := mgo.Dial("localhost")
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Create indices
	// 인덱스 값으로 email 을 사용하며, 그 값은 고유하다
	if err = db.Copy().DB("st_more").C("users").EnsureIndex(mgo.Index{
		Key:    []string{"email"},
		Unique: true,
	}); err != nil {
		log.Fatal(err)
	}

	//---------------
	// Route & Server
	//---------------

	// Initialize handler
	h := &handler.Handler{DB: db}

	// Route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Something!\n")
	})                                  // 인덱스
	e.POST("/sign-up/", h.SignUpNormal) // 회원 가입
	e.POST("/admin/", h.SignUpAdmin)    // 관리자 회원 가입
	e.GET("/activate/", h.Activate)     // 이메일 회원 활성화
	e.POST("/sign-in/", h.SignIn)       // 로그인

	e.POST("/story/", h.CreateStory) // 스토리 생성
	e.GET("/story/", h.ListStory) // 스토리 리스트
	e.GET("/story/count/", h.CountStory) // 스토리 총 갯수
	e.GET("/story/:story_id", h.RetrieveStory) // 스토리 디테일
	e.PATCH("/story/:story_id", h.PatchStory) // 스토리 수정
	e.DELETE("/story/:story_id", h.DestroyStory) // 스토리 삭제

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
