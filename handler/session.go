package handler

import (
	// Default package
	"os"
	"io"
	"fmt"
	"strings"
	"strconv"
	"math/rand"
	"path/filepath"
	"mime/multipart"
	// Third Party package
	"github.com/labstack/echo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	// User package
	"github.com/backend/model"
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

func (h *Handler) FindPost(c echo.Context, p *model.Post, q string) (err error) {

	// Get IDs
	postID := c.Param(fmt.Sprintf("%s_id", q))

	// Find story in database
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C(q).
		FindId(bson.ObjectIdHex(postID)).
		One(p); err != nil {
		if err == mgo.ErrNotFound {
			return echo.ErrNotFound
		}
		return
	}
	return
}

func (h *Handler) MapAuthorNickname(c echo.Context, p *model.Post) (err error) {
	// 포스트 객체에 담긴 userID 를 이용해 AuthorNickname 을 구하는 함수
	u := new(model.User)
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C(USER).
		FindId(p.AuthorID).One(u); err != nil {
		if err == mgo.ErrNotFound {
			// 유저를 찾을 수 없는 경우 닉네임에 "탈퇴한 회원" 값을 줌
			p.AuthorNickname = "탈퇴한 회원"
		}
		return
	}
	p.AuthorNickname = u.Nickname

	return
}

func (h *Handler) UploadThumbnail(c echo.Context, s *model.Post, file *multipart.FileHeader) (err error) {
	// File open
	src, err := file.Open()
	defer src.Close()
	if err != nil {
		return
	}

	// 파일명 해쉬값 생성
	letters := []byte ("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var sliceStr []string // 빈 문자 타입 슬라이스 선언
	for i := 0; i < 6; i++ { // 6차례 순회
		randInt := rand.Intn(len(letters))        // letters 슬라이스의 길이 범위 안에서 랜덤 정수 생성
		convertedStr := string(letters[randInt])  // 랜덤 정수를 인덱스로 하는 letters 의 요소를 꺼내 문자열로 변환
		sliceStr = append(sliceStr, convertedStr) // 빈 슬라이스에 랜덤 문자열 추가
	}
	resultStr := strings.Join(sliceStr, "")                 // 랜덤 문자열을 하나의 문자열로 합침
	resultInt := strconv.FormatInt(rand.Int63n(100000), 10) // 100000 이하의 정수를 랜덤으로 생성해 문자열로 변환

	// 파일명 생성하기
	baseFileName := file.Filename                           // 원본 파일 이름
	extension := filepath.Ext(baseFileName)                 // 확장자 추출
	realName := strings.TrimSuffix(baseFileName, extension) // 확장자 제외한 이름 추출
	newFileName := realName + "_" + resultStr + resultInt + extension // 새로운 파일명 생성

	// 정적 파일 루트 디렉터리 생성하기: 디렉터리가 없으면 생성할 것
	temporaryPath, _ := filepath.Abs("../bin/")
	assetsPath := filepath.Join(temporaryPath, "/assets/")
	if _, err := os.Stat(assetsPath); os.IsNotExist(err) {
		os.Mkdir(assetsPath, 0777)
	}

	// 정적 파일 아래 필진 디렉터리 생성하기
	authorPathValue := "/" + s.AuthorID.Hex() + "/" // authorID 로 된 path 값 생성
	authorPath := filepath.Join(assetsPath, authorPathValue)
	if _, err := os.Stat(authorPath); os.IsNotExist(err) {
		os.Mkdir(authorPath, 0777)
	}

	// 파일 생성하기
	dst, err := os.Create(filepath.Join(authorPath, newFileName)) // 파일 이름에 랜덤 문자열 추가하여 저장
	defer dst.Close()
	if err != nil {
		return
	}

	// 복사하기
	if _, err = io.Copy(dst, src); err != nil {
		return
	}

	// 파일 주소명을 Story object 에 넣기
	thumbnailURL := c.Scheme() + "://" + c.Request().Host + "/assets" + authorPathValue + newFileName
	s.Thumbnail = thumbnailURL
	return
}

func (h *Handler) PatchThumbnail(s *model.Post) (err error) {
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(DBName).C(STORY).
		Update(
		bson.M{"_id": s.ID},
		bson.M{"$set":
		bson.M{
			"thumbnail": s.Thumbnail}}); err != nil {
		return
	}
	return
}
