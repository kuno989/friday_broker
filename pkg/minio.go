package pkg

import (
	"bytes"
	"context"
	"github.com/google/wire"
	miniogo "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

var (
	DefaultMinioConfig = MinioConfig{
		URI:       "localhost:9000",
		Bucket:    "malwares",
	}
	MinioProviderSet = wire.NewSet(NewMinio, ProvideMinioConfig)
)

type MinioConfig struct {
	URI       string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"accessKey"`
	SecretKey string `mapstructure:"secretKey"`
	Bucket    string `mapstructure:"bucket"`
}

func ProvideMinioConfig(cfg *viper.Viper) (MinioConfig, error) {
	minio := DefaultMinioConfig
	err := cfg.UnmarshalKey("minio", &minio)
	return minio, err
}

type Minio struct {
	Config MinioConfig
	Client *miniogo.Client
}

func NewMinio(cfg MinioConfig) (*Minio, error) {
	client, err := miniogo.New(cfg.URI, &miniogo.Options{
		Creds: credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
	})
	if err != nil {
		return nil, err
	}
	return &Minio{
		Config: cfg,
		Client: client,
	}, nil
}

func (m Minio) Upload(ctx context.Context, fileInfo *multipart.FileHeader) (miniogo.UploadInfo, error) {
	file, err := fileInfo.Open()
	if err != nil {
		return miniogo.UploadInfo{}, err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	contentType := http.DetectContentType(content)
	buf := bytes.NewBuffer(content)
	info, err := m.Client.PutObject(ctx, m.Config.Bucket, "sample/"+fileInfo.Filename, buf, fileInfo.Size, miniogo.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return miniogo.UploadInfo{}, err
	}
	return info, nil
}

func (m Minio) Download(ctx context.Context, key string, file io.Writer) error {
	reader, err := m.Client.GetObject(ctx, m.Config.Bucket, key, miniogo.GetObjectOptions{})
	if err != nil {
		return err
	}
	stat, err := reader.Stat()
	if err != nil {
		return err
	}
	_, err = io.CopyN(file, reader, stat.Size)
	if err != nil {
		return err
	}
	return nil
}
