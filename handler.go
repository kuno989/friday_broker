package connect

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/kuno989/friday_connect/connect/schema"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)
const (
	regex = "^[a-f0-9]{64}$"
	queued       = iota
	processing   = iota
	finished     = iota
	vmProcessing = iota
)

func (s *Server) index(c echo.Context) error {
	return c.JSON(http.StatusOK, schema.Response{
		Message: "Success",
		Description: "Friday Broker 정상 작동 중",
	})
}

func (s *Server) JobStart(c echo.Context) error {
	logrus.Info("파일 작업 시작")
	sha256 := strings.ToLower(c.Param("sha256"))
	matched, _ := regexp.MatchString(regex, sha256)
	if !matched {
		return c.JSON(http.StatusBadRequest, schema.FileResponse{
			Message:     "sha256 포멧이 아닙니다",
			Description: "잘못된 hash 요청입니다 : " + sha256,
		})
	}
	b, err := ioutil.ReadAll(c.Request().Body)
	ctx := context.Background()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	file, err := s.ms.FileSearch(ctx, sha256)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(http.StatusOK, schema.FileResponse{
				Sha256:      sha256,
				Message:     err.Error(),
				Description: "파일을 찾을 수 없습니다",
			})
		}
	}
	err = json.Unmarshal(b, &file)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(http.StatusInternalServerError, schema.FileResponse{
				Sha256:      sha256,
				Message:     err.Error(),
				Description: "작업 중 에러가 발생하였습니다",
			})
		}
	}

	file.Status = vmProcessing
	file, err = s.ms.FileUpdate(ctx, file)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(http.StatusInternalServerError, schema.FileResponse{
				Sha256:      sha256,
				Message:     err.Error(),
				Description: "업데이트 중 에러가 발생하였습니다",
			})
		}
	}


	return c.JSON(http.StatusOK, schema.Response{
		Message: "Success",
		Description: "작업 상태 변경 완료",
	})
}