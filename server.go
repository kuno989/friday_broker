package connect

import (
	"bytes"
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
		if err := vm.Start(); err != nil {
			logrus.Errorf("machine start failure %s", err)
		}
		logrus.Infof("%s sandbox start", vm.Name)
		s.vmRequest(resp.MinioObjectKey, resp.Sha256, resp.FileType, "download")
	} else {
		logrus.Infof("job won't start because file is not pe")
	}

	//s.vmRequest(resp.MinioObjectKey, resp.Sha256, "download")

	//if err := machine.Save(); err != nil {
	//	logrus.Errorf("machine save failure %s", err)
	//}
	//logrus.Infof("%s sandbox saved", machine.Name)

	return nil
}
