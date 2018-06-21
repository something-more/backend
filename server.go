package main

import (
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

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//CORS WhiteList
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderAccept,
			echo.HeaderContentType,
			echo.HeaderXRequestedWith,
			"X-XSRF-TOKEN",
		},
		AllowCredentials: true,
	}))
	// XSRF Token
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "header:X-XSRF-TOKEN",
		ContextKey: "csrftoken",
		CookieName: "csrftoken",
	}))

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

	// Initialize handler
	h := &handler.Handler{DB: db}

	// Route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Something!\n")
	}) // 인덱스
	e.POST("/signup/", h.SignUp) // 회원 가입

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
