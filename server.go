package connect

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/wire"
	"github.com/kuno989/friday_connect/connect/pkg"
	"github.com/kuno989/friday_connect/connect/schema"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"github.com/terra-farm/go-virtualbox"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var (
	DefaultServerConfig = ServerConfig{
		Debug: true,
	}
	ServerProviderSet = wire.NewSet(NewServer, ProvideServerConfig)
)

type ServerConfig struct {
	Debug          bool     `mapstructure:"debug"`
	AgentURI       string   `mapstructure:"agent_uri"`
	AgentPort      string   `mapstructure:"agent_port"`
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	MaxFileSize    int64    `mapstructure:"maxFileSize"`
	VBoxManage     string   `mapstructure:"vbox"`
	VBoxName       string   `mapstructure:"vbox_name"`
	VBoxSnapshot   string   `mapstructure:"vbox_snapshot"`
	VBoxDumpPath   string   `mapstructure:"vm_dump_path"`
}

func ProvideServerConfig(cfg *viper.Viper) (ServerConfig, error) {
	sc := DefaultServerConfig
	err := cfg.Unmarshal(&sc)
	return sc, err
}

type Server struct {
	*echo.Echo
	Config ServerConfig
	ms     *pkg.Mongo
	Rb     *pkg.RabbitMq
	minio  *pkg.Minio
}

func NewServer(cfg ServerConfig, ms *pkg.Mongo, rb *pkg.RabbitMq, minio *pkg.Minio) *Server {
	s := &Server{
		Echo:   echo.New(),
		Config: cfg,
		ms:     ms,
		Rb:     rb,
		minio:  minio,
	}
	s.HideBanner = true
	s.HidePort = true
	var allowedOrigins []string
	if cfg.Debug {
		allowedOrigins = append(cfg.AllowedOrigins, "http://localhost:3000", "*")
	} else {
		if len(cfg.AllowedOrigins) == 0 {
			allowedOrigins = []string{"http://localhost:3000"}
		} else {
			allowedOrigins = cfg.AllowedOrigins
		}
	}
	s.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization", "Access-Control-allow-Methods", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
	}))
	s.RegisterHandlers()
	return s
}

// virtualbox 와 통신하기 위한 rest api
func (s *Server) RegisterHandlers() {
	api := s.Group("/api")
	api.GET("/", s.index)
	api.PUT("/jobStart/:sha256", s.JobStart)
	api.POST("/jobEnd/:sha256", s.JobEnd)
}

func (s *Server) AmqpHandler(msg amqp.Delivery) error {
	body := bytes.ReplaceAll(msg.Body, []byte("NaN"), []byte("0"))
	var resp schema.ResponseObject
	if err := json.Unmarshal(body, &resp); err != nil {
		return errors.New("ailed to parse message body")
		if err := msg.Reject(false); err != nil {
			logrus.Errorf("failed to reject message %v", err)
		}
	}
	if resp.FileType == "pe"{
		logrus.Infof("job start")
		vm, err := virtualbox.GetMachine("win7")
		if err != nil {
			logrus.Errorf("can not find machine %s", err)
		}
		logrus.Infof("%s sandbox found", vm.Name)
		logrus.Infof("cpu %v, memory %v", vm.CPUs, vm.Memory)

		logrus.Infof("VM 상태 : %s",vm.State)
		if vm.State == "poweroff"{
			s.VmRestore()
			for {
				time.Sleep(1 * time.Second)
				if vm.State == "poweroff" || vm.State == "saved"{
					vm.State = "running"
					break
				}
			}
		}
		if vm.State == "running" || vm.State == "saved"{
			for {
				time.Sleep(1 * time.Second)
				s.VmStart()
				logrus.Infof("%s sandbox start", vm.Name)
				if vm.State == "running" || vm.State == "saved"{
					vm.State = "running"
					break
				}
			}
		}
		time.Sleep(5 * time.Second)
		s.VmResolution()
		time.Sleep(10 * time.Second)
		logrus.Info("분석 대기중")
		if vm.State == "running"{
			s.vmRequest(resp.MinioObjectKey, resp.Sha256, resp.FileType, "download", "POST")
		}
	} else {
		ctx := context.Background()
		file, err := s.ms.FileSearch(ctx, resp.Sha256)
		if err != nil {
			if !errors.Is(err, mongo.ErrNoDocuments) {
				logrus.Errorf("DB에서 파일을 찾을 수 없습니다", err)
			}
		}
		file.IsNotPE = true
		file.Status = 5
		file, err = s.ms.FileUpdate(ctx, file)
		if err != nil {
			logrus.Errorf("DB 업데이트 중 에러가 발생하였습니다", err)
		}

		logrus.Infof("job won't start because file is not pe")
	}
	return nil
}
