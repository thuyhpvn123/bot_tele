package services

import (
	"bytes"
	"fmt"
	// "math/big"
	"net/http"
	// "strconv"
	// "time"
	"github.com/meta-node-blockchain/meta-node/cmd/usdtnoti/internal/database/models"
	// "github.com/meta-node-blockchain/meta-node/cmd/usdtnoti/internal/utils"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
)
const (
	pin_emoji                string = "📌"
	annoucement_emoji        string = "📣"
	Title_Mint_By_Controller string = "\\[MINT BY CONTROLLER\\]"
)
type TeleService struct {
	chatID   string
	botToken string
}

func NewTeleService(chatId string, botToken string) *TeleService {
	return &TeleService{
		chatID:   chatId,
		botToken: botToken,
	}
}

func (s *TeleService) SendNoti(msg []byte) error {
	fmt.Println("msg la:",msg)
	jsonStr := []byte(
		fmt.Sprintf(`{"chat_id": "%v", "text": "%v"}`, s.chatID, string(msg)),
	)
	resp, err := http.Post(
		fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.botToken),
		"application/json",
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		logger.Debug("resp ", resp)
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (s *TeleService) SendMintNoti(mintInfo *models.MintInfo) error {
	// humanTime := time.Unix(int64(subInfo.Time), 0).Format("2006-01-02 15:04:05")
	// normalAmount := new(big.Int).Div(subInfo.Amount, big.NewInt(1000000)).String()
	var buffer bytes.Buffer
	buffer.WriteString("📣📣📣[Mint Noti - USDT]📣📣📣\n")
	buffer.WriteString(fmt.Sprintf("%s%s%s*%s*%s%s%s\n \n", annoucement_emoji, annoucement_emoji, annoucement_emoji, Title_Mint_By_Controller, annoucement_emoji, annoucement_emoji, annoucement_emoji))
	buffer.WriteString(fmt.Sprintf("%sController: _%s_ \n", pin_emoji, mintInfo.ControllerAddress))
	buffer.WriteString(fmt.Sprintf("%sRecipient: _%s_ \n", pin_emoji, mintInfo.RecipientAddress))
	buffer.WriteString(fmt.Sprintf("%sAmount: *%s* \n", pin_emoji, mintInfo.Amount))
	buffer.WriteString(fmt.Sprintf("%sTotalMint: *%s* \n", pin_emoji, mintInfo.TotalMint))
	return s.SendNoti(buffer.Bytes())
}

