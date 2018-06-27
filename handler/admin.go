package handler

import (
	// Default package
	"net/http"
	// Third Party package
	"github.com/labstack/echo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	// User package
	"github.com/backend/model"
)

func (h *Handler) ListUsers(c echo.Context) (err error) {

	// Validate user
	u := new(model.User)

	userID := userIDFromToken(c)

	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("st_more").C("users").
		FindId(bson.ObjectIdHex(userID)).
		One(u); err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "존재하지 않는 계정입니다",
			}
		}
		return
	}

	// Validate admin
	if isAdminFromToken(c) == false {
		return &echo.HTTPError{
			Code: http.StatusUnauthorized,
			Message: "이 계정은 관리자가 아닙니다",
		}
	}

	// Find users
	var users []*model.User
	if err = db.DB("st_more").C("users").
		Find(nil).
		Sort("-is_admin").
		Sort("-is_staff").
		All(&users); err != nil {
		return
	}

	return c.JSON(http.StatusOK, users)
}