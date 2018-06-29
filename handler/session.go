package handler

import (
	"github.com/labstack/echo"
	"github.com/globalsign/mgo"
	"github.com/backend/model"
)

func (h *Handler) FindUser(c echo.Context, u *model.User) (err error) {
	db := h.DB.Clone()
	defer db.Close()

	if err = db.DB("st_more").C("users").FindId(u.ID).One(u); err != nil {
		if err == mgo.ErrNotFound {
			return echo.ErrNotFound
		}
		return
	}
	return
}
