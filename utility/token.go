package utility

import (
	// Default package
	"time"
	// Third-party package
	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
	// User package
	"github.com/backend/model"

)

func CreateJWT(u *model.User) *jwt.Token {
	// Create token
	// HS256 알고리즘으로 인코딩
	// 단방향 암호화 알고리즘인 RS256과 달리 양방향 암호화 알고리즘이므로 디코딩이 가능함
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	// 유저 정보를 담는다
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = u.ID
	claims["email"] = u.Email
	claims["nickname"] = u.Nickname
	claims["isActive"] = u.IsActive
	claims["isStaff"] = u.IsStaff
	claims["isAdmin"] = u.IsAdmin
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // 토큰 유효시간: 72시간

	return token
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

func UserNicknameFromToken(c echo.Context) string {
	// JWT 를 통해 닉네임을 체크하는 헬퍼 함수
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["nickname"].(string)
}

func IsAdminFromToken(c echo.Context) bool {
	// JWT 를 통해 관리자 여부를 체크하는 헬퍼 함수
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["isAdmin"].(bool)
}
