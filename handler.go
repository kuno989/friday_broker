package connect

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/kuno989/friday_connect/connect/schema"
	models "github.com/kuno989/friday_connect/connect/schema/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/terra-farm/go-virtualbox"
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
	vmFinished   = iota
)

func (s *Server) index(c echo.Context) error {
	return c.JSON(http.StatusOK, schema.Response{
		Message: "Success",
		Description: "Friday Broker 정상 작동 중",
	})
}

func (s *Server) JobEnd(c echo.Context) error {
	logrus.Info("작업 종료")
	sha256 := strings.ToLower(c.Param("sha256"))
	matched, _ := regexp.MatchString(regex, sha256)
	if !matched {
		return c.JSON(http.StatusBadRequest, schema.FileResponse{
			Message:     "sha256 포멧이 아닙니다",
			Description: "잘못된 hash 요청입니다 : " + sha256,
		})
	}

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var resp = models.DBModel{}
	err = json.Unmarshal(b, &resp)

	ctx := context.Background()

	var mongoResp = schema.File{}
	buff, err := json.Marshal(mongoResp)

	file, err := s.ms.FileSearch(ctx, sha256)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			logrus.Errorf("DB에서 파일을 찾을 수 없습니다", err)
		}
	}

	err = json.Unmarshal(buff, &file)
	if err != nil {
		logrus.Errorf("mongodb %v", err)
	}
	file.Status = vmFinished
	file.DBModel = resp

	var pid string
	//var childPID string
	//var malName string

	for _, mal := range file.DBModel.ProcessCreate {
		if strings.Contains(mal.ProcessPath, sha256) && !strings.Contains(mal.Operation, `C:\Windows\system32\cmd.exe`){
			pid = mal.ChildPID
		}
	}

	if pid == ""{
		for _, mal := range file.DBModel.ProcessCreate {
			if strings.Contains(mal.Operation, sha256){
				pid = mal.ChildPID
			}
		}
	}

	file.DBModel.MalPid = pid
	file.DBModel.MalName = sha256 + ".exe"

	for _, sub := range file.DBModel.ProcessCreate{
		if strings.Contains(sub.ProcessName, "svchost"){
			if strings.Contains(sub.Operation, "DllHost"){
				for _, cf := range file.DBModel.CreateFile{
					if cf.PID == sub.ChildPID{
						file.DBModel.SubPid = sub.ChildPID
					}
				}
			}
		}
	}

	file.IsNotPE = false

	file, err = s.ms.FileUpdate(ctx, file)
	if err != nil {
		logrus.Errorf("DB 업데이트 중 에러가 발생하였습니다", err)
	}
	logrus.Info("DB 업데이트 완료")

	logrus.Info("메모리 덤프 중..")
	dumpPath := s.VmDump(sha256)
	_, err = s.minio.DumpUpload(ctx, dumpPath)
	if err != nil {
		logrus.Errorf("메모리 덤프 업로드 중 에러 %v", err)
	}
	logrus.Info("메모리 덤프 업로드 완료")

	logrus.Info("VM 종료")
	vm, err := virtualbox.GetMachine("win7")
	if err != nil {
		logrus.Errorf("can not find machine %v", err)
	}
	vm.Poweroff()


	return c.JSON(http.StatusOK, schema.Response{
		Message: "Success",
		Description: "작업 완료",
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

	var resp = schema.ResponseAgent{}
	err = json.Unmarshal(b, &resp)

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
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(http.StatusInternalServerError, schema.FileResponse{
				Sha256:      sha256,
				Message:     err.Error(),
				Description: "업데이트 중 에러가 발생하였습니다",
			})
		}
	}

	s.vmRequest(resp.MinioObjectKey, resp.Sha256, resp.FileType, "start", "POST")

	file.Status = vmProcessing
	_, err = s.ms.FileUpdate(ctx, file)
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
