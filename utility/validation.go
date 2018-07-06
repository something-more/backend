package utility

import (
	// Default package
	"net/http"
	// Third-party package
	"github.com/labstack/echo"
)

func EmptyValueValidation(c echo.Context) (err error) {

	if c.FormValue("title") == "" || c.FormValue("content") == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "제목이나 내용을 반드시 입력해야 합니다",
		}
		return
	}
	return
}

func AdminValidation(c echo.Context) (err error) {

	if IsAdminFromToken(c) == false {
		return echo.ErrUnauthorized
	}
	return
}
