package handler

import (
	"github.com/labstack/echo"
	"github.com/backend/utility"
	"github.com/globalsign/mgo/bson"
	"github.com/backend/model"
	"net/http"
	"strconv"
)

func (h *Handler) CreateNotice(c echo.Context) (err error) {
	// Find user in database
	userID := utility.UserIDFromToken(c)
	if err = h.FindUser(userID); err != nil {
		return
	}

	// Validate admin
	if err = utility.AdminValidation(c); err != nil {
		return
	}

	// Bind notice object
	n := &model.Post{
		ID:     bson.NewObjectId(),
		Author: userID, // 저자를 표시하기 위해 u.ID 를 삽입
	}

	if err = c.Bind(n); err != nil {
		return
	}

	// Empty Value Validation
	if err = utility.EmptyValueValidation(c); err != nil {
		return
	}

	// Add FormValue in Post Instance
	n.Title = c.FormValue("title")
	n.Content = c.FormValue("content")
	n.DateCreated = c.FormValue("date_created")
	n.DateModified = ""
	n.IsPublished = true

	// Save Post
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C(NOTICE).Insert(n); err != nil {
		return
	}

	return c.JSON(http.StatusCreated, n)
}

func (h *Handler) ListNotice(c echo.Context) (err error) {
	// Get query params
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	// Default pagination
	// 페이지 당 최대 20개의 글만 쿼리
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 20
	}

	// List stories from database
	userID := utility.UserIDFromToken(c)
	var notices []*model.Post

	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C(NOTICE).
		Find(bson.M{"author": userID}).
		Select(bson.M{"content": 0}). // 내용은 받아오지 않음으로써 응답시간 단축
		Sort("-date_created"). // 생성일자 역순으로 정렬
		Skip((page - 1) * limit).
		Limit(limit).
		All(&notices); err != nil {
		return
	}

	return c.JSON(http.StatusOK, notices)
}