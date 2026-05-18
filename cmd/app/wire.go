//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"github.com/makesalekz/stores/internal/biz"
	"github.com/makesalekz/stores/internal/conf"
	"github.com/makesalekz/stores/internal/data"
	"github.com/makesalekz/stores/internal/server"
	"github.com/makesalekz/stores/internal/service"
)

func wireApp(*conf.Bootstrap, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
