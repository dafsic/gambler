package bet

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dafsic/gambler/config"
	"github.com/dafsic/gambler/lib/mylog"
	"github.com/dafsic/gambler/modules/channels"
	"github.com/dafsic/gambler/utils"
	"go.uber.org/fx"
)

type Bet interface {
	Working()
	Stop()
}

type BetImpl struct {
	pk   []byte
	pool string
	rc   chan interface{}
	wc   chan interface{}
	qc   chan bool
	l    *utils.Logger
}

func NewBet(lc fx.Lifecycle, cfg config.ConfigI, log mylog.Logging, chanMgr channels.ChanManager) Bet {
	pk := cfg.GetElem("betpk").(string)
	b := &BetImpl{
		pk:   []byte(pk),
		pool: cfg.GetElem("pooladdress").(string),
		rc:   chanMgr.GetChan("bet"),
		wc:   chanMgr.GetChan("bethash"),
		l:    log.GetLogger("bet"),
		qc:   make(chan bool, 1),
	}

	lc.Append(fx.Hook{
		// app.start调用
		OnStart: func(ctx context.Context) error {
			// 这里不能阻塞
			go b.Working()
			return nil
		},
		// app.stop调用，收到中断信号的时候调用app.stop
		OnStop: func(ctx context.Context) error {
			go b.Stop()
			return nil
		},
	})

	return b
}

var BetModule = fx.Options(fx.Provide(NewBet))

func (b *BetImpl) Working() {
	for {
		select {
		case amount := <-b.rc:
			hash, err := b.betting(amount.(int64))
			if err != nil {
				b.l.Error(err.Error())
				break
			}
			b.wc <- hash
		case <-b.qc:
			return
		}
	}
	// for amount := range b.rc {
	// 	hash, err := b.betting(amount.(int64))
	// 	if err != nil {
	// 		b.l.Error(err.Error())
	// 		continue
	// 	}
	// 	b.wc <- hash
	// }
}

func (b *BetImpl) Stop() {
	b.qc <- true
}

func (b *BetImpl) betting(amount int64) (string, error) {
	url := "http://127.0.0.1:8090/wallet/easytransferbyprivate"

	// decode, _, err := base58.CheckDecode(to)
	// if err != nil {
	// 	return "", err
	// }
	// hexTo := fmt.Sprintf("41%x", decode)
	payload := strings.NewReader(fmt.Sprintf("{\"privateKey\":\"%s\",\"toAddress\":\"%s\",\"amount\":%d}", string(b.pk), b.pool, amount))
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w%s", err, utils.LineNo())
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	// TODO: 解析出交易hash
	var r TxResult

	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", fmt.Errorf("%w%s", err, utils.LineNo())
	}
	return r.T.H, nil
}

type Transaction struct {
	H string `json:"txId"`
}

type TxResult struct {
	T Transaction `json:"transaction"`
}
