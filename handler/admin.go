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
	"github.com/backend/utility"
)

func (h *Handler) ValidateAdmin(c echo.Context) (err error) {

	u := new(model.User)

	userID := utility.UserIDFromToken(c)

	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C("users").
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

	if utility.IsAdminFromToken(c) == false {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "이 계정은 관리자가 아닙니다",
		}
	}

	return
}

func (h *Handler) ListUsers(c echo.Context) (err error) {

	// Validate Admin
	if err = h.ValidateAdmin(c); err != nil {
		return
	}

	// Find users
	var users []*model.User
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C("users").
		Find(nil).
		Select(bson.M{"password": 0}).
		Sort("-is_admin").
		Sort("-is_staff").
		All(&users); err != nil {
		return
	}

	return c.JSON(http.StatusOK, users)
}

func (h *Handler) UpdateUserAuth(c echo.Context) (err error) {

	// Validate Admin
	if err = h.ValidateAdmin(c); err != nil {
		return
	}

	// Bind object
	u := new(model.User)
	if err = c.Bind(u); err != nil {
		return
	}

	userEmail := c.Param("user_email")

	// Update user authentication
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C("users").
		Update(
		bson.M{"email": userEmail},
		bson.M{"$set":
		bson.M{
			"is_admin": u.IsAdmin,
			"is_staff": u.IsStaff}}); err != nil {
		return
	}

	return c.NoContent(http.StatusOK)
}
