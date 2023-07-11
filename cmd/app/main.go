package main

import (
	"github.com/Tyz3/nymgraph/cmd/app/config"
	"github.com/Tyz3/nymgraph/internal/repository"
	"github.com/Tyz3/nymgraph/internal/service"
	"github.com/Tyz3/nymgraph/internal/state"
	"github.com/Tyz3/nymgraph/internal/view"
	"github.com/Tyz3/nymgraph/pkg/client/sqlite3"
	"github.com/Tyz3/nymgraph/pkg/utils"
)

func init() {
	if err := utils.SaveResource(config.CfgFileName, config.CfgBin); err != nil {
		panic(err)
	}
	if err := utils.SaveResource(config.DataDBFileName, config.DataDBBin); err != nil {
		panic(err)
	}
}

func main() {
	// Конфигурация
	cfg, err := config.NewConfig(config.CfgFileName)
	if err != nil {
		panic(err)
	}

	// Текущие состояния приложения
	states := state.NewState(cfg)

	// Клиент подключения к БД
	client, err := sqlite3.NewClient(sqlite3.Config{DBFileName: "data.db"})
	if err != nil {
		panic(err)
	}

	// Инициализация репозитория
	repo := repository.NewRepository(client)

	// Инициализация контроллеров(сервисов)
	serv := service.NewService(repo, states)

	// Инициализация и запуск приложения
	app := view.NewApp("Nymgraph", serv)
	app.Run()
	app.Close()
}
