package pkg

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

var (
	DefaultRabbitMqConfig = RabbitMqConfig{
		URI:           "amqp://guest:guest@message-broker:5672/",
	}
	RabbitmqProviderSet = wire.NewSet(NewRabbitMq, ProvideRabbitMqConfig)
)

type RabbitMqConfig struct {
	URI           string `mapstructure:"uri"`
	VmQueue       string `mapstructure:"vmQueue"`
}

func ProvideRabbitMqConfig(cfg *viper.Viper) (RabbitMqConfig, error) {
	rb := DefaultRabbitMqConfig
	err := cfg.UnmarshalKey("rabbitmq", &rb)
	return rb, err
}

type RabbitMq struct {
	Config RabbitMqConfig
	Client *amqp.Connection
}

func NewRabbitMq(cfg RabbitMqConfig) (*RabbitMq, func(), error) {
	client, err := amqp.Dial(cfg.URI)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		client.Close()
	}
	return &RabbitMq{
		Config: cfg,
		Client: client,
	}, cleanup, nil
}

func (r *RabbitMq) Channel() (*amqp.Channel, error) {
	ch, err := r.Client.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}
