package utility

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/dgrijalva/jwt-go"
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

func UserIDFromToken(c echo.Context) string {
	// 다른 메서드 안에서 JWT 를 통해 DB 상의 ID 를 꺼내오는 헬퍼 함수
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["id"].(string)
}

func UserEmailFromToken(c echo.Context) string {
	// JWT 를 통해 이메일을 체크하는 헬퍼 함수
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["email"].(string)
}

func IsAdminFromToken(c echo.Context) bool {
	// JWT 를 통해 관리자 여부를 체크하는 헬퍼 함수
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["isAdmin"].(bool)
}
