package main

import (
	"github.com/craftdome/nymgraph/cmd/app/config"
	"github.com/craftdome/nymgraph/internal/repository"
	"github.com/craftdome/nymgraph/internal/service"
	"github.com/craftdome/nymgraph/internal/state"
	"github.com/craftdome/nymgraph/internal/view"
	"github.com/craftdome/nymgraph/pkg/client/sqlite3"
	"github.com/craftdome/nymgraph/pkg/utils"
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
