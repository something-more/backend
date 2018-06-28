package handler

import (
	"github.com/labstack/echo"
	"github.com/backend/model"
	"github.com/globalsign/mgo/bson"
	"net/http"
	"github.com/globalsign/mgo"
)

func (h *Handler) CreateBoard(c echo.Context) (err error) {
	// Bind object
	u := &model.User{
		ID: bson.ObjectIdHex(userIDFromToken(c)),
	}
	b := &model.Board{
		ID: bson.NewObjectId(),
		Author: u.Email, // 저자를 표시하기 위해 u.ID 를 삽입
	}

	if err = c.Bind(b); err != nil {
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

	// Add FormValue in Board Instance
	b.Title = c.FormValue("title")
	b.Content = c.FormValue("content")
	b.DateCreated = c.FormValue("date_created")
	b.DateModified = ""

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
	if err = db.DB("st_more").C("board").Insert(b); err != nil {
		return
	}
	return c.JSON(http.StatusCreated, b)
}

func (h *Handler) listBoard(c echo.Context) (err error) {

	var boards []*model.Board

	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("st_more").C("board").
		Find(nil).
		Sort("-date_created").
		All(&boards); err != nil {
			return
	}

	return c.JSON(http.StatusOK, boards)
}