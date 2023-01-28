package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dafsic/gambler/config"
	"github.com/dafsic/gambler/lib/mylog"
	"github.com/dafsic/gambler/modules"
	"github.com/dafsic/gambler/modules/pkmanager"
	"github.com/dafsic/gambler/utils"
	"go.uber.org/fx"
)

type TrxClient interface {
	Transfer(from, to string, amount int64, token string) (string, error)
	TransferTrx(from, to string, amount int64) (string, error)
	TransferContract(from, to string, amount int64, contract string) (string, error)
	GetBalance(addr string) (int64, error)
}

type TrxClientImpl struct {
	Node  string
	pkMgr pkmanager.PKManager
	l     *utils.Logger
}

func NewTrxClient(cfg config.ConfigI, log mylog.Logging, pkmgr pkmanager.PKManager) TrxClient {
	t := &TrxClientImpl{
		Node:  cfg.GetElem("trxnode").(string),
		pkMgr: pkmgr,
		l:     log.GetLogger("client"),
	}
	t.l.Info("Init...")

	return t
}

func (t *TrxClientImpl) TransferTrx(from, to string, amount int64) (string, error) {
	url := fmt.Sprintf("http://%s/wallet/easytransferbyprivate", t.Node)
	// decode, _, err := base58.CheckDecode(to)
	// if err != nil {
	// 	return "", err
	// }
	// hexTo := fmt.Sprintf("41%x", decode)
	pk, ok := t.pkMgr.GetPrivateKey(from)
	if !ok {
		return "", fmt.Errorf("%s%s", "can not get private key", utils.LineNo())
	}
	payload := strings.NewReader(fmt.Sprintf("{\"privateKey\":\"%s\",\"toAddress\":\"%s\",\"amount\":%d}", pk, to, amount))
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

// TransferContract 通过contract地址，使用pk作为私钥，向to地址发送amount的token
func (t *TrxClientImpl) TransferContract(from, to string, amount int64, contract string) (string, error) {

	msg, err := t.createTx(from, to, contract, amount)
	if err != nil {
		return "", fmt.Errorf("%w%s", err, utils.LineNo())
	}

	pk, ok := t.pkMgr.GetPrivateKey(from)
	if !ok {
		return "", fmt.Errorf("can not get pk of:%s,%s", from, utils.LineNo())
	}

	signMsg, err := t.sign(msg, pk)
	if err != nil {
		return "", fmt.Errorf("%w%s", err, utils.LineNo())
	}

	rst, err := t.broadcast(signMsg)
	if err != nil {
		return "", fmt.Errorf("%w%s", err, utils.LineNo())
	}

	var r struct {
		Txid string `json:"txid"`
	}
	err = json.Unmarshal(rst, &r)
	if err != nil {
		return "", fmt.Errorf("%w%s", err, utils.LineNo())
	}
	if r.Txid == "" {
		return "", fmt.Errorf("%s%s", rst, utils.LineNo())
	}
	return r.Txid, nil
}

func (t *TrxClientImpl) createTx(from, to, contract string, amount int64) ([]byte, error) {
	url := fmt.Sprintf("http://%s/wallet/triggersmartcontract", t.Node)

	para := fmt.Sprintf("%064s%064x", to[2:], amount)
	payload := strings.NewReader(fmt.Sprintf("{\"contract_address\":\"%s\",\"function_selector\":\"transfer(address,uint256)\",\"parameter\":\"%s\",\"fee_limit\":100000000,\"call_value\":0,\"owner_address\":\"%s\"}", contract, para, from))
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("%w%s", err, utils.LineNo())
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("%w%s", err, utils.LineNo())
	}

	//解析出交易hash
	var r struct {
		Transaction json.RawMessage `json:"transaction"`
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return []byte{}, fmt.Errorf("%w%s", err, utils.LineNo())
	}

	return r.Transaction, nil
}

func (t *TrxClientImpl) sign(msg []byte, pk string) ([]byte, error) {
	url := fmt.Sprintf("http://%s/wallet/gettransactionsign", t.Node)
	payload := strings.NewReader(fmt.Sprintf("{\"transaction\":%s,\"privateKey\":\"%s\"}", msg, pk))
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("%w%s", err, utils.LineNo())
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("%w%s", err, utils.LineNo())
	}

	return body, nil
}

func (t *TrxClientImpl) broadcast(msg []byte) ([]byte, error) {
	url := fmt.Sprintf("http://%s/wallet/broadcasttransaction", t.Node)

	req, _ := http.NewRequest("POST", url, bytes.NewReader(msg))
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("%w%s", err, utils.LineNo())
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("%w%s", err, utils.LineNo())
	}

	return body, nil
}

func (t *TrxClientImpl) Transfer(from, to string, amount int64, token string) (string, error) {
	switch strings.ToUpper(token) {
	case "TRX":
		return t.TransferTrx(from, to, amount)
	case "USDT":
		return t.TransferContract(from, to, amount, modules.USDT_CONTRACT)
	default:
		return "", fmt.Errorf("not support token,%s", utils.LineNo())
	}
}

var ClientModule = fx.Options(fx.Provide(NewTrxClient))
