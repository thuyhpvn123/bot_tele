package app

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/meta-node-blockchain/meta-node/cmd/client"
	"github.com/meta-node-blockchain/meta-node/cmd/usdtnoti/internal/config"
	// "github.com/meta-node-blockchain/meta-node/cmd/usdtnoti/internal/database"
	"github.com/meta-node-blockchain/meta-node/cmd/usdtnoti/internal/handler"
	"github.com/meta-node-blockchain/meta-node/cmd/usdtnoti/internal/services"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"github.com/meta-node-blockchain/meta-node/types"

	c_config "github.com/meta-node-blockchain/meta-node/cmd/client/pkg/config"
)

func PreflightHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

type App struct {
	Config      *config.AppConfig
	ChainClient *client.Client
	EventChan   chan types.EventLogs
	StopChan    chan bool

	UsdtHandler *handler.UsdtHandler
}

func NewApp(
	configFilePath string,
	logLevel int,
) (*App, error) {
	var loggerConfig = &logger.LoggerConfig{
		Flag:    logLevel,
		Outputs: []*os.File{os.Stdout},
	}
	logger.SetConfig(loggerConfig)

	app := &App{}
	// load config
	var err error
	app.Config, err = config.LoadConfig(configFilePath)
	if err != nil {
		logger.Error(fmt.Sprintf("error when loading config %v", err))
		return nil, err
	}

	
	// event channel
	app.ChainClient, err = client.NewClient(
		&c_config.ClientConfig{
			Version_:                app.Config.MetaNodeVersion,
			PrivateKey_:             app.Config.PrivateKey_,
			ParentAddress:           app.Config.NodeAddress,
			ParentConnectionAddress: app.Config.NodeConnectionAddress,
			DnsLink_:                app.Config.DnsLink(),
		},
	)

	if err != nil {
		logger.Error(fmt.Sprintf("error when create chain client %v", err))
		return nil, err
	}

	app.EventChan, err = app.ChainClient.Subcribes(
		common.HexToAddress(app.Config.StorageAddress),
		[]common.Address{
			common.HexToAddress(app.Config.UsdtAddress),
		},
	)
	if err != nil {
		logger.Error(fmt.Sprintf("error when create chain client %v", err))
		return nil, err
	}

	// create card abi
	reader, err := os.Open(app.Config.UsdtABIPath) // * Unit Test
	if err != nil {
		logger.Error("Error occured while read create card smart contract abi")
		return nil, err
	}
	defer reader.Close()

	usdtAbi, err := abi.JSON(reader)
	if err != nil {
		logger.Error("Error occured while parse create card smart contract abi")
		return nil, err
	}

	app.UsdtHandler = handler.NewUsdtHandler(
		app.ChainClient,
		common.HexToAddress(
			app.Config.UsdtAddress,
		),
		&usdtAbi,
		app.Config.MintHash,		
		services.NewTeleService(app.Config.ChatID, app.Config.BotToken),
	)

	// Initialize the Gin router
	r := gin.Default()
	// Initialize cors config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowCredentials = true

	r.Use(cors.New(corsConfig))
	// Initialize services
	go func() {
		err = r.Run(app.Config.APIAddress)
		if err != nil {
			log.Fatal("Failed to start the server")
		}
	}()

	return app, nil
}

func (app *App) Run() {
	app.StopChan = make(chan bool)
	for {
		select {
		case <-app.StopChan:
			return
		case eventLogs := <-app.EventChan:
			logger.Debug(eventLogs)
			app.UsdtHandler.HandleEvent(eventLogs)
		}
	}
}

func (app *App) Stop() error {
	app.ChainClient.Close()
	// defer database.CloseDB()

	logger.Warn("App Stopped")
	return nil
}
