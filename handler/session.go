package handler

import (
	"github.com/labstack/echo"
	"github.com/globalsign/mgo"
	"github.com/backend/model"
	"github.com/globalsign/mgo/bson"
)

func (h *Handler) FindUser(id string) (err error) {
	db := h.DB.Clone()
	defer db.Close()

	if err = db.DB("st_more").C("users").FindId(bson.ObjectIdHex(id)).One(nil); err != nil {
		if err == mgo.ErrNotFound {
			return echo.ErrNotFound
		}
		return
	}
	return
}

func (h *Handler) FindStory(c echo.Context, s *model.Post) (err error) {

	// Get IDs
	storyID := c.Param("story_id")

	// Find story in database
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("st_more").C("stories").
		Find(bson.M{"_id": bson.ObjectIdHex(storyID)}).
		One(s); err != nil {
		if err == mgo.ErrNotFound {
			return echo.ErrNotFound
		}
		return
	}
	return
}

func (h *Handler) FindBoard(c echo.Context, b *model.Board) (err error) {

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