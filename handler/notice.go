package handler

import (
	"github.com/labstack/echo"
	"github.com/backend/utility"
	"github.com/globalsign/mgo/bson"
	"github.com/backend/model"
	"net/http"
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
	if err = db.DB(DBName).C(STORY).Insert(n); err != nil {
		return
	}

	return c.JSON(http.StatusCreated, n)
}
