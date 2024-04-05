package handler

import (
	// "encoding/json"
	"math/big"
	// "slices"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	e_common "github.com/ethereum/go-ethereum/common"
	"github.com/meta-node-blockchain/meta-node/cmd/client"
	"github.com/meta-node-blockchain/meta-node/cmd/usdtnoti/internal/database/models"
	"github.com/meta-node-blockchain/meta-node/cmd/usdtnoti/internal/services"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"github.com/meta-node-blockchain/meta-node/types"
)

type UsdtHandler struct {
	chainClient     *client.Client
	retailSCAddress e_common.Address
	usdtSCAbi   *abi.ABI
	mintHash      string
	teleServ        *services.TeleService
}

func NewUsdtHandler(
	chainClient *client.Client,
	retailSCAddress e_common.Address,
	usdtSCAbi *abi.ABI,
	mintHash string,
	teleServ *services.TeleService,
) *UsdtHandler {
	return &UsdtHandler{
		chainClient:     chainClient,
		retailSCAddress: retailSCAddress,
		usdtSCAbi:   usdtSCAbi,
		mintHash:      mintHash,
		teleServ:        teleServ,
	}
}


func (h *UsdtHandler) HandleEvent(
	events types.EventLogs,
) {
	fmt.Println("11111111111111",events)
	for _, v := range events.EventLogList() {
		fmt.Println("v.Topics():",v.Topics())
		switch v.Topics()[0] {
		case h.mintHash:
			{
				eventResult := make(map[string]interface{})
				err := h.usdtSCAbi.UnpackIntoMap(eventResult, "MintByController", e_common.FromHex(v.Data()))
				if err != nil {
					logger.Error("error when unpack into map")
					continue
				}
				// typ := eventResult["typ"].(string)
				addController := eventResult["_controller"].(common.Address)
				addRecipient := eventResult["_recipient"].(common.Address)
				amount := eventResult["_amount"].(*big.Int)
				totalMint := eventResult["_totalMint"].(*big.Int)

				mintHistory := &models.MintInfo{
					ControllerAddress:      addController.String(),
					RecipientAddress:       addRecipient.String(),
					Amount:          		amount,
					TotalMint:				totalMint,
				}
				h.teleServ.SendMintNoti(mintHistory)
			}
		
		default:
			continue
		}
	}
}
