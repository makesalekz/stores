//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"gitlab.calendaria.team/services/stores/internal/biz"
	"gitlab.calendaria.team/services/stores/internal/conf"
	"gitlab.calendaria.team/services/stores/internal/data"
	"gitlab.calendaria.team/services/stores/internal/server"
	"gitlab.calendaria.team/services/stores/internal/service"
)

func wireApp(*conf.Bootstrap, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
