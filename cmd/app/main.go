package main

import (
	"fmt"
	"github.com/Tyz3/nymgraph/cmd/app/config"
	"github.com/Tyz3/nymgraph/internal/controller"
	"github.com/Tyz3/nymgraph/internal/repository"
	"github.com/Tyz3/nymgraph/internal/state"
	"github.com/Tyz3/nymgraph/internal/view"
	"github.com/Tyz3/nymgraph/pkg/client/sqlite3"
	"github.com/Tyz3/nymgraph/pkg/logger"
	"github.com/Tyz3/nymgraph/pkg/utils"
	"os"
)

func init() {
	if err := utils.SaveResource(config.CfgFileName, config.CfgBin); err != nil {
		logger.Log.ERROR.Printf("%s\n", err)
		return
	}
	if err := utils.SaveResource(config.DataDBFileName, config.DataDBBin); err != nil {
		logger.Log.ERROR.Printf("%s\n", err)
		return
	}
}

func main() {
	// Конфигурация
	cfg, err := config.NewConfig(config.CfgFileName)
	if err != nil {
		logger.Log.ERROR.Printf("%s\n", err)
		return
	}

	// Текущие состояния приложения
	states := state.NewState(cfg)

	// Клиент подключения к БД
	client, err := sqlite3.NewClient(sqlite3.Config{DBFileName: "data.db"})
	if err != nil {
		logger.Log.ERROR.Printf("%s\n", err)
		return
	}

	// Инициализация репозитория
	repo := repository.NewRepository(client)

	// Инициализация контроллеров(сервисов)
	controllers := controller.NewController(repo, states)

	//controllers.NymClient.AddSelfAddressHandler(func(r *response.SelfAddress) {
	//	states.SelfAddress = r.Address
	//	logger.Log.INFO.Printf("SelfAddress: %s\n", states.SelfAddress)
	//})
	//
	//controllers.NymClient.AddErrorHandler(func(r *response.Error) {
	//	logger.Log.ERROR.Printf("nym-client: %s\n", r.Message)
	//})
	//
	//controllers.NymClient.AddReceivedHandler(func(r *response.Received) {
	//	logger.Log.INFO.Printf("Received: %s SenderTag: %s\n", r.Message, r.SenderTag)
	//})
	//
	//controllers.NymClient.AddLaneQueueLengthHandler(func(r *response.LaneQueueLength) {
	//
	//})
	//
	//if err := controllers.NymClient.Dial(states.GetNymClientUrl()); err != nil {
	//	logger.Log.ERROR.Printf("%s\n", err)
	//	return
	//}
	//
	//go func() {
	//	if err := controllers.NymClient.ListenAndServe(); err != nil && err != we.ErrClosed {
	//		logger.Log.ERROR.Printf("%s\n", err)
	//	}
	//}()
	//
	//getSelfAddrReq := nym.NewGetSelfAddress()
	//if err := controllers.NymClient.SendRequestAsText(getSelfAddrReq); err != nil {
	//	logger.Log.ERROR.Printf("%s\n", err)
	//}
	//
	//receive := nym.NewSend("msg1", "CpE6HPVeWammXnTB7Au2AoChJmMRXtFEA7hgXxSMA4Bi.rQZrjMeaacNNez2Bzyt6eEpjz1LxZ3TX6khhTAh3RcN@Fo4f4SQLdoyoGkFae5TpVhRVoXCF8UiypLVGtGjujVPf")
	//if err := controllers.NymClient.SendRequestAsText(receive); err != nil {
	//	logger.Log.ERROR.Printf("%s\n", err)
	//}

	//anon := nym.NewSendAnonymous("message1", "CpE6HPVeWammXnTB7Au2AoChJmMRXtFEA7hgXxSMA4Bi.rQZrjMeaacNNez2Bzyt6eEpjz1LxZ3TX6khhTAh3RcN@Fo4f4SQLdoyoGkFae5TpVhRVoXCF8UiypLVGtGjujVPf", 0)
	//if err := controllers.NymClient.SendRequestAsText(anon); err != nil {
	//	logger.Log.ERROR.Printf("%s\n", err)
	//}
	//
	//reply := nym.NewReply("reply1", "LagaXrsPeuRUAcH2cx8xo4")
	//if err := controllers.NymClient.SendRequestAsText(reply); err != nil {
	//	logger.Log.ERROR.Printf("%s\n", err)
	//}

	// Инициализация и запуск приложения
	app := view.NewApp("Nymgraph", controllers)
	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
