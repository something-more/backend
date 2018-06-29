package handler

import (
	"github.com/labstack/echo"
	"github.com/globalsign/mgo"
	"github.com/backend/model"
	"github.com/globalsign/mgo/bson"
	"fmt"
)

const DBName = "st_more"
const USER = "users"
const STORY = "story"
const BOARD = "board"
const NOTICE = "notice"

func (h *Handler) FindUser(id string) (err error) {
	db := h.DB.Clone()
	defer db.Close()

	if err = db.DB(DBName).C(USER).FindId(bson.ObjectIdHex(id)).One(nil); err != nil {
		if err == mgo.ErrNotFound {
			return echo.ErrNotFound
		}
		return
	}
	return
}

func (h *Handler) FindPost(c echo.Context, s *model.Post, q string) (err error) {

	// Get IDs
	postID := c.Param(fmt.Sprintf("%s_id", q))

	// Find story in database
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C(q).
		Find(bson.M{"_id": bson.ObjectIdHex(postID)}).
		One(s); err != nil {
		if err == mgo.ErrNotFound {
			return echo.ErrNotFound
		}
		return
	}
	return
}
