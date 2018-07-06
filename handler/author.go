package handler

import (
	// Default package
	"fmt"
	"strconv"
	"net/http"
	// Third Party package
	"github.com/labstack/echo"
	"github.com/globalsign/mgo/bson"
	// User package
	"github.com/backend/model"
)

func (h *Handler) ListAuthors(c echo.Context) (err error) {
	// Find users
	var users []*model.User
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C(USER).
		Find(bson.M{"is_staff": true}).
		Select(bson.M{"password": 0}).
		Sort("-is_admin").
		Sort("-is_staff").
		All(&users); err != nil {
		return
	}

	return c.JSON(http.StatusOK, users)
}

func (h *Handler) ListStoryAuthor(c echo.Context) (err error) {
	// Get query params
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	// Get Author IDs
	AuthorID := c.Param(fmt.Sprint("author_id"))

	// Default pagination
	// 페이지 당 최대 15개의 글만 쿼리
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 15
	}

	// Find story in database
	var stories []*model.Post
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C(STORY).
		Find(bson.M{
		"author_id":    bson.ObjectIdHex(AuthorID),
		"is_published": true}).
		Select(bson.M{"content": 0}).
		Sort("-date_created").
		Skip((page - 1) * limit).
		Limit(limit).
		All(&stories); err != nil {
		return
	}

	// stories 슬라이스 순회
	for _, story := range stories {
		h.MapAuthorNickname(c, story)
	}

	return c.JSON(http.StatusOK, stories)
}

func (h *Handler) CountStoryAuthor(c echo.Context) (err error) {
	// Get Author IDs
	AuthorID := c.Param(fmt.Sprint("author_id"))

	// int type 변수 지정
	var count int

	db := h.DB.Clone()
	defer db.Close()
	if count, err = db.DB(DBName).C(STORY).
		Find(bson.M{
		"author_id":    bson.ObjectIdHex(AuthorID),
		"is_published": true}).
		Count(); err != nil {
		return
	}

	// int type count 를 ascii 로 변환해서 리턴
	return c.String(http.StatusOK, strconv.Itoa(count))
}
