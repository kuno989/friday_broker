//go:build wireinject
// +build wireinject

package connect

import (
	"context"
	"github.com/kuno989/friday_connect/connect/pkg"

	"github.com/google/wire"
	"github.com/spf13/viper"

)

func InitializeServer(ctx context.Context, cfg *viper.Viper) (*Server, func(), error) {
	panic(wire.Build(ServerProviderSet, pkg.MongoProviderSet, pkg.RabbitmqProviderSet, pkg.MinioProviderSet))
}
