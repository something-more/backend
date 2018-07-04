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
	baseFileName := file.Filename                                     // 원본 파일 이름
	extension := filepath.Ext(baseFileName)                           // 확장자 추출
	realName := strings.TrimSuffix(baseFileName, extension)           // 확장자 제외한 이름 추출
	newFileName := realName + "_" + resultStr + resultInt + extension // 새로운 파일명 생성

	// 디렉터리 생성하기
	filePath, err := filepath.Abs("../bin/assets/") // 디렉터리가 없으면 생성할 것
	if err != nil {
		os.MkdirAll(filePath, 0777)
	}

	// 파일 생성하기
	dst, err := os.Create(filepath.Join(filePath, newFileName)) // 파일 이름에 랜덤 문자열 추가하여 저장
	defer dst.Close()
	if err != nil {
		return
	}

	// 복사하기
	if _, err = io.Copy(dst, src); err != nil {
		return
	}

	// 파일 주소명을 Story object 에 넣기
	thumbnailURL := c.Scheme() + "://" + c.Request().Host + "/assets/" + newFileName
	s.Thumbnail = thumbnailURL
	return
}
