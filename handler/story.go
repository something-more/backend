package handler

import (
	// Default package
	"strconv"
	"net/http"
	// Third Party package
	"github.com/labstack/echo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	// User package
	"github.com/backend/model"
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
	s.DateCreated = c.FormValue("date_created")
	s.DateModified = ""
	s.IsPublished = false

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

func (h *Handler) ListStory(c echo.Context) (err error) {
	userID := userIDFromToken(c)
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

	// Retrieve posts from database
	stories := []*model.Story{}
	db := h.DB.Clone()
	if err = db.DB("st_more").C("stories").
		Find(bson.M{"author": userID}).
		Sort("-date_created"). // 생성일자 역순으로 정렬
		Skip((page - 1) * limit).
		Limit(limit).
		All(&stories); err != nil {
		return
	}
	defer db.Close()

	return c.JSON(http.StatusOK, stories)
}

func (h *Handler) CountStory(c echo.Context) (err error) {
	userID := userIDFromToken(c)

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

func (h *Handler) GetStory(c echo.Context, s *model.Story) (err error) {

	// Get IDs
	userID := userIDFromToken(c)
	storyID := c.Param("story_id")

	// Find story in database
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("st_more").C("stories").
		Find(bson.M{"author":userID, "_id": bson.ObjectIdHex(storyID)}).
		One(s); err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "스토리를 찾을 수 없습니다",
			}
		}
		return
	}

	return
}

func (h *Handler) RetrieveStory(c echo.Context) (err error) {
	// Object bind
	s := new(model.Story)
	if err = c.Bind(s); err != nil {
		return
	}

	// Find story in database
	if err = h.GetStory(c, s); err != nil {
		return
	}

	return c.JSON(http.StatusOK, s)
}

func (h *Handler) PatchStory(c echo.Context) (err error) {
	// Object bind
	s := new(model.Story)
	if err = c.Bind(s); err != nil {
		return
	}

	// Find story in database
	if err = h.GetStory(c, s); err != nil {
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
	if err = h.GetStory(c, s); err != nil {
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