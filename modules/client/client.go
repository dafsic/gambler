package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dafsic/gambler/config"
	"github.com/dafsic/gambler/lib/mylog"
	"github.com/dafsic/gambler/utils"
	"go.uber.org/fx"
)

type TrxClient interface {
	Betting(amount int64, pk, pool string) (string, error)
	GetBalance(addr string) (int64, error)
}

type TrxClientImpl struct {
	Node string
	l    *utils.Logger
}

func NewTrxClient(cfg config.ConfigI, log mylog.Logging) TrxClient {
	t := &TrxClientImpl{
		Node: cfg.GetElem("trxnode").(string),
		l:    log.GetLogger("client"),
	}
	t.l.Info("Init...")

	return t
}

func (t *TrxClientImpl) Betting(amount int64, pk, pool string) (string, error) {
	url := fmt.Sprintf("http://%s/wallet/easytransferbyprivate", t.Node)
	// decode, _, err := base58.CheckDecode(to)
	// if err != nil {
	// 	return "", err
	// }
	// hexTo := fmt.Sprintf("41%x", decode)
	payload := strings.NewReader(fmt.Sprintf("{\"privateKey\":\"%s\",\"toAddress\":\"%s\",\"amount\":%d}", pk, pool, amount))
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w%s", err, utils.LineNo())
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	//解析出交易hash
	var r TxResult
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", fmt.Errorf("%w%s", err, utils.LineNo())
	}
	if r.T.H == "" {
		return "", fmt.Errorf("%s%s", body, utils.LineNo())
	}
	return r.T.H, nil
}

type Transaction struct {
	H string `json:"txId"`
}

type TxResult struct {
	T Transaction `json:"transaction"`
}

func (t *TrxClientImpl) GetBalance(addr string) (int64, error) {
	url := fmt.Sprintf("http://%s/wallet/getaccount", t.Node)

	payload := strings.NewReader(fmt.Sprintf("{\"address\":\"%s\"}", addr))
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("%w%s", err, utils.LineNo())
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var r struct {
		Balance int64 `jsom:"balance"`
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		return 0, fmt.Errorf("%w%s", err, utils.LineNo())
	}

	return r.Balance, nil
}

var ClientModule = fx.Options(fx.Provide(NewTrxClient))
