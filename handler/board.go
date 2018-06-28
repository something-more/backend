package handler

import (
	"github.com/labstack/echo"
	"github.com/backend/model"
	"github.com/globalsign/mgo/bson"
	"net/http"
	"github.com/globalsign/mgo"
	"strconv"
)

func (h *Handler) CreateBoard(c echo.Context) (err error) {
	// Bind object
	u := &model.User{
		ID: bson.ObjectIdHex(userIDFromToken(c)),
		Email: userEmailFromToken(c),
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

func (h *Handler) ListBoard(c echo.Context) (err error) {

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	// Default pagination
	// 페이지 당 최대 20개의 글만 쿼리
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 15
	}

	var boards []*model.Board

	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("st_more").C("board").
		Find(nil).
		Sort("-date_created").
		Skip((page - 1) * limit).
		Limit(limit).
		All(&boards); err != nil {
			return
	}

	return c.JSON(http.StatusOK, boards)
}

func (h *Handler) CountBoard(c echo.Context) (err error) {

	// int type 변수 지정
	var count int

	// Get count of stories from database
	db := h.DB.Clone()
	defer db.Close()
	if count, err = db.DB("st_more").C("board").
		Find(nil).
		Count(); err != nil {
		return
	}

	// int type count 를 ascii 로 변환해서 리턴
	return c.String(http.StatusOK, strconv.Itoa(count))
}

func (h *Handler) GetBoard(c echo.Context, b *model.Board) (err error) {

	// Get IDs
	boardID := c.Param("board_id")

	// Find board in database
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("st_more").C("board").
		Find(bson.M{"_id": bson.ObjectIdHex(boardID)}).
		One(b); err != nil {
		if err == mgo.ErrNotFound {
			return echo.ErrNotFound
		}
		return
	}
	return
}

func (h *Handler) RetrieveBoard(c echo.Context) (err error) {
	// Object bind
	b := new(model.Board)
	if err = c.Bind(b); err != nil {
		return
	}

	// Find story in database
	if err = h.GetBoard(c, b); err != nil {
		return
	}

	return c.JSON(http.StatusOK, b)
}

func (h *Handler) PatchBoard(c echo.Context) (err error) {
	// Object bind
	b := new(model.Board)
	if err = c.Bind(b); err != nil {
		return
	}

	// Find user
	userEmail := userEmailFromToken(c)

	// Find story in database
	if err = h.GetBoard(c, b); err != nil {
		return
	}

	// Add FormValues in Story Instance
	b.Title = c.FormValue("title")
	b.Content = c.FormValue("content")
	b.DateModified = c.FormValue("date_modified")

	// Update story in database
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("st_more").C("board").
		Update(
		bson.M{"_id": b.ID, "author": userEmail},
		bson.M{"$set":
		bson.M{
			"title":         b.Title,
			"content":       b.Content,
			"date_modified": b.DateModified}}); err != nil {
		return
	}

	return c.JSON(http.StatusOK, b)
}

func (h *Handler) DestroyBoard(c echo.Context) (err error) {
	// Object bind
	b := new(model.Board)
	if err = c.Bind(b); err != nil {
		return
	}

	// Find user
	userEmail := userEmailFromToken(c)

	// Find story in database
	if err = h.GetBoard(c, b); err != nil {
		return
	}

	// Destroy board in database
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("st_more").C("board").
		Remove(bson.M{"_id": b.ID, "author": userEmail}); err != nil {
			return
	}

	return c.NoContent(http.StatusNoContent)
}