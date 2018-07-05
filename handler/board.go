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

func (h *Handler) CreateBoard(c echo.Context) (err error) {
	// Find user in database
	userID := utility.UserIDFromToken(c)
	if err = h.FindUser(userID); err != nil {
		return
	}

	// Bind board object
	b := &model.Post{
		ID:       bson.NewObjectId(),
		AuthorID: bson.ObjectIdHex(userID), // 저자를 표시하기 위해 userID 삽입
	}

	if err = c.Bind(b); err != nil {
		return
	}

	// Empty Value Validation
	if err = utility.EmptyValueValidation(c); err != nil {
		return
	}

	// Add FormValue in Post Instance
	b.Title = c.FormValue("title")
	b.Content = c.FormValue("content")
	b.DateCreated = c.FormValue("date_created")
	b.DateModified = ""
	b.IsPublished = true

	// Save Post
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C(BOARD).Insert(b); err != nil {
		return
	}

	return c.JSON(http.StatusCreated, b)
}

func (h *Handler) ListBoard(c echo.Context) (err error) {
	// Get query params
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

	// List boards from database
	var boards []*model.Post

	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C(BOARD).
		Find(nil).
		Select(bson.M{"content": 0}). // 내용은 받아오지 않음으로써 응답시간 단축
		Sort("-date_created").
		Skip((page - 1) * limit).
		Limit(limit).
		All(&boards); err != nil {
		return
	}

	// Boards 슬라이스 순회하며 닉네임 매핑
	for _, board := range boards {
		h.MapAuthorNickname(c, board)
	}

	return c.JSON(http.StatusOK, boards)
}

func (h *Handler) CountBoard(c echo.Context) (err error) {

	// int type 변수 지정
	var count int

	// Get count of stories from database
	db := h.DB.Clone()
	defer db.Close()
	if count, err = db.DB(DBName).C(BOARD).
		Find(nil).
		Count(); err != nil {
		return
	}

	// int type count 를 ascii 로 변환해서 리턴
	return c.String(http.StatusOK, strconv.Itoa(count))
}

func (h *Handler) RetrieveBoard(c echo.Context) (err error) {
	// Object bind
	b := new(model.Post)
	if err = c.Bind(b); err != nil {
		return
	}

	// Find story in database
	if err = h.FindPost(c, b, BOARD); err != nil {
		return
	}

	// Map AuthorNickname
	h.MapAuthorNickname(c, b)

	return c.JSON(http.StatusOK, b)
}

func (h *Handler) PatchBoard(c echo.Context) (err error) {
	// Find user in database
	userID := utility.UserIDFromToken(c)
	if err = h.FindUser(userID); err != nil {
		return
	}

	// Object bind
	b := new(model.Post)
	if err = c.Bind(b); err != nil {
		return
	}

	// Find story in database
	if err = h.FindPost(c, b, BOARD); err != nil {
		return
	}

	// Add FormValues in Post Instance
	b.Title = c.FormValue("title")
	b.Content = c.FormValue("content")
	b.DateModified = c.FormValue("date_modified")

	// Update story in database
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C(BOARD).
		Update(
		bson.M{"_id": b.ID},
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
	// Find user in database
	userID := utility.UserIDFromToken(c)
	if err = h.FindUser(userID); err != nil {
		return
	}

	// Object bind
	b := new(model.Post)
	if err = c.Bind(b); err != nil {
		return
	}

	// Find story in database
	if err = h.FindPost(c, b, BOARD); err != nil {
		return
	}

	// Destroy board in database
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C(BOARD).
		Remove(bson.M{"_id": b.ID}); err != nil {
		return
	}

	return c.NoContent(http.StatusNoContent)
}
