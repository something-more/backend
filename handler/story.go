package handler

import (
	// Default package
	"strconv"
	"net/http"
	// Third Party package
	"github.com/labstack/echo"
	"github.com/globalsign/mgo/bson"
	// User package
	"github.com/backend/model"
	"github.com/backend/utility"
)

func (h *Handler) CreateStory(c echo.Context) (err error) {
	// Find user in database
	userID := utility.UserIDFromToken(c)
	if err = h.FindUser(userID); err != nil {
		return
	}

	// Bind story object
	s := &model.Story{
		ID:     bson.NewObjectId(),
		Author: userID, // 저자를 표시하기 위해 u.ID 를 삽입
	}

	if err = c.Bind(s); err != nil {
		return
	}

	// Empty Value Validation
	if err = utility.EmptyValueValidation(c); err != nil {
		return
	}

	// Add FormValue in Story Instance
	s.Title = c.FormValue("title")
	s.Content = c.FormValue("content")
	s.DateCreated = c.FormValue("date_created")
	s.DateModified = ""
	s.IsPublished = false

	// Save Story
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("st_more").C("stories").Insert(s); err != nil {
		return
	}

	return c.JSON(http.StatusCreated, s)
}

func (h *Handler) ListStory(c echo.Context) (err error) {
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
	var stories []*model.Story

	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("st_more").C("stories").
		Find(bson.M{"author": userID}).
		Sort("-date_created"). // 생성일자 역순으로 정렬
		Skip((page - 1) * limit).
		Limit(limit).
		All(&stories); err != nil {
		return
	}

	return c.JSON(http.StatusOK, stories)
}

func (h *Handler) CountStory(c echo.Context) (err error) {
	userID := utility.UserIDFromToken(c)

	// int type 변수 지정
	var count int

	// Get count of stories from database
	db := h.DB.Clone()
	defer db.Close()
	if count, err = db.DB("st_more").C("stories").
		Find(bson.M{"author": userID}).
		Count(); err != nil {
		return
	}

	// int type count 를 ascii 로 변환해서 리턴
	return c.String(http.StatusOK, strconv.Itoa(count))
}

func (h *Handler) RetrieveStory(c echo.Context) (err error) {
	// Object bind
	s := new(model.Story)
	if err = c.Bind(s); err != nil {
		return
	}

	// Find story in database
	if err = h.FindStory(c, s); err != nil {
		return
	}

	return c.JSON(http.StatusOK, s)
}

func (h *Handler) PatchStory(c echo.Context) (err error) {
	// Find user in database
	userID := utility.UserIDFromToken(c)
	if err = h.FindUser(userID); err != nil {
		return
	}

	// Object bind
	s := new(model.Story)
	if err = c.Bind(s); err != nil {
		return
	}

	// Find story in database
	if err = h.FindStory(c, s); err != nil {
		return
	}

	// Add FormValues in Story Instance
	s.Title = c.FormValue("title")
	s.Content = c.FormValue("content")
	s.DateModified = c.FormValue("date_modified")
	s.IsPublished, _ = strconv.ParseBool(c.FormValue("is_published"))

	// Update story in database
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("st_more").C("stories").
		Update(
		bson.M{"_id": s.ID},
		bson.M{"$set":
		bson.M{
			"title":         s.Title,
			"content":       s.Content,
			"date_modified": s.DateModified,
			"is_published":  s.IsPublished}}); err != nil {
		return
	}

	return c.JSON(http.StatusOK, s)
}

func (h *Handler) DestroyStory(c echo.Context) (err error) {
	// Object bind
	s := new(model.Story)
	if err = c.Bind(s); err != nil {
		return
	}

	// Find story in database
	if err = h.FindStory(c, s); err != nil {
		return
	}

	// Destroy story in database
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("st_more").C("stories").
		Remove(bson.M{"_id": s.ID}); err != nil {
		return
	}

	return c.NoContent(http.StatusNoContent)
}
