package handler

import (
	"github.com/labstack/echo"
	"github.com/backend/model"
	"github.com/globalsign/mgo/bson"
	"net/http"
	"github.com/globalsign/mgo"
)

func (h * Handler) CreateStory(c echo.Context) (err error) {
	// Object bind
	// 유저는 JWT 에서 알아낸 DB 상의 ID를 16진수 디코딩을 하여 찾아낸다
	u := &model.User{
		ID: bson.ObjectIdHex(userIDFromToken(c)),
	}
	s := &model.Story{
		ID: bson.NewObjectId(),
		Author: u.ID.Hex(), // 저자를 표시하기 위해 u.ID 를 삽입
	}
	if err = c.Bind(s); err != nil {
		return
	}

	// Validation
	if c.FormValue("title") == "" || c.FormValue("content") == "" {
		return &echo.HTTPError{
			Code: http.StatusBadRequest,
			Message: "제목이나 내용을 반드시 입력해야 합니다",
		}
		return
	}

	// Add FormValue in Story Instance
	s.Title = c.FormValue("title")
	s.Content = c.FormValue("content")

	// Find user
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("st_more").C("users").FindId(u.ID).One(u); err != nil {
		if err == mgo.ErrNotFound {
			return echo.ErrNotFound
		}
		return
	}

	// Save Story
	if err = db.DB("st_more").C("stories").Insert(s); err != nil {
		return
	}
	return c.JSON(http.StatusCreated, s)
}