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
	s := &model.Post{
		ID:       bson.NewObjectId(),
		AuthorID: bson.ObjectIdHex(userID), // 저자를 표시하기 위해 userID 삽입
	}

	if err = c.Bind(s); err != nil {
		return
	}

	// Empty Value Validation
	if err = utility.EmptyValueValidation(c); err != nil {
		return
	}

	// Thumbnail Upload Validation
	file, err := c.FormFile("thumbnail")
	// 썸네일이 입력되어 에러가 발생되지 않을 때에만 업로드 썸네일 함수 실행
	if err == nil {
		if err = h.UploadThumbnail(c, s, file); err != nil {
			return
		}
	}

	// Add FormValue in Post Instance
	s.Title = c.FormValue("title")
	s.Content = c.FormValue("content")
	s.DateCreated = c.FormValue("date_created")
	s.Category = c.FormValue("category")
	s.DateModified = ""
	s.IsPublished = false

	// Save Post
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C(STORY).Insert(s); err != nil {
		return
	}

	return c.JSON(http.StatusCreated, s)
}

func (h *Handler) ListStory(c echo.Context) (err error) {
	// Get query params
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	// Default pagination
	// 페이지 당 최대 15개의 글만 쿼리
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 15
	}

	// List stories from database
	userID := utility.UserIDFromToken(c)
	var stories []*model.Post

	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C(STORY).
		Find(bson.M{"author_id": bson.ObjectIdHex(userID)}).
		Select(bson.M{"content": 0}). // 내용은 받아오지 않음으로써 응답시간 단축
		Sort("-date_created"). // 생성일자 역순으로 정렬
		Skip((page - 1) * limit).
		Limit(limit).
		All(&stories); err != nil {
		return
	}

	return c.JSON(http.StatusOK, stories)
}

func (h *Handler) ClientListStory(c echo.Context) (err error) {
	// Get query params
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	// Default pagination
	// 페이지 당 최대 4개의 글만 쿼리
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 4
	}

	// List stories from database
	var stories []*model.Post

	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C(STORY).
		Find(bson.M{"is_published": true}). // 발행 승인된 게시글만 쿼리
		Select(bson.M{"content": 0}). // 내용은 받아오지 않음으로써 응답시간 단축
		Sort("-date_created"). // 생성일자 역순으로 정렬
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

func (h *Handler) CountStory(c echo.Context) (err error) {

	// int type 변수 지정
	var count int

	// Get count of stories from database
	userID := utility.UserIDFromToken(c)

	db := h.DB.Clone()
	defer db.Close()
	if count, err = db.DB(DBName).C(STORY).
		Find(bson.M{"author_id": bson.ObjectIdHex(userID)}).
		Count(); err != nil {
		return
	}

	// int type count 를 ascii 로 변환해서 리턴
	return c.String(http.StatusOK, strconv.Itoa(count))
}

func (h *Handler) RetrieveStory(c echo.Context) (err error) {
	// Object bind
	s := new(model.Post)
	if err = c.Bind(s); err != nil {
		return
	}

	// Find story in database
	if err = h.FindPost(c, s, STORY); err != nil {
		return
	}

	// Map AuthorNickname
	h.MapAuthorNickname(c, s)

	return c.JSON(http.StatusOK, s)
}

func (h *Handler) PatchStory(c echo.Context) (err error) {
	// Find user in database
	userID := utility.UserIDFromToken(c)
	if err = h.FindUser(userID); err != nil {
		return
	}

	// Object bind
	s := new(model.Post)
	if err = c.Bind(s); err != nil {
		return
	}

	// Find story in database
	if err = h.FindPost(c, s, STORY); err != nil {
		return
	}

	// Thumbnail Upload Validation
	file, err := c.FormFile("thumbnail")
	// 썸네일이 입력되어 에러가 발생되지 않을 때에만 업로드 썸네일 함수 실행
	if err == nil {
		if err = h.UploadThumbnail(c, s, file); err != nil {
			return
		}
		if err = h.PatchThumbnail(s); err != nil {
			return
		}
	}

	// Add FormValues in Post Instance
	s.Title = c.FormValue("title")
	s.Content = c.FormValue("content")
	s.DateModified = c.FormValue("date_modified")
	s.IsPublished, _ = strconv.ParseBool(c.FormValue("is_published"))
	s.Category = c.FormValue("category")

	// Update story in database
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C(STORY).
		Update(
		bson.M{"_id": s.ID},
		bson.M{"$set":
		bson.M{
			"title":         s.Title,
			"content":       s.Content,
			"date_modified": s.DateModified,
			"is_published":  s.IsPublished,
			"category":      s.Category}}); err != nil {
		return
	}

	return c.JSON(http.StatusOK, s)
}

func (h *Handler) ChangePublishStory(c echo.Context) (err error) {
	// Find user in database
	userID := utility.UserIDFromToken(c)
	if err = h.FindUser(userID); err != nil {
		return
	}

	// Object bind
	s := new(model.Post)
	if err = c.Bind(s); err != nil {
		return
	}

	// Find story in database
	if err = h.FindPost(c, s, STORY); err != nil {
		return
	}

	// Add FormValues in Post Instance
	s.DateModified = c.FormValue("date_modified")
	s.IsPublished, _ = strconv.ParseBool(c.FormValue("is_published"))

	// Update story in database
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C(STORY).
		Update(
		bson.M{"_id": s.ID},
		bson.M{"$set":
		bson.M{
			"date_modified": s.DateModified,
			"is_published":  s.IsPublished}}); err != nil {
		return
	}

	return c.JSON(http.StatusOK, s)
}

func (h *Handler) DestroyStory(c echo.Context) (err error) {
	// Object bind
	s := new(model.Post)
	if err = c.Bind(s); err != nil {
		return
	}

	// Find story in database
	if err = h.FindPost(c, s, STORY); err != nil {
		return
	}

	// Destroy story in database
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C(STORY).
		Remove(bson.M{"_id": s.ID}); err != nil {
		return
	}

	return c.NoContent(http.StatusNoContent)
}
