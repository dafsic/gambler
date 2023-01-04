package transferhandler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dafsic/gambler/lib/mylog"
	"github.com/dafsic/gambler/modules/client"
	"github.com/dafsic/gambler/modules/server"
	"github.com/dafsic/gambler/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Agent interface {
	RegisterRouters()
	SwitchOn(*gin.Context)
	SwitchOff(*gin.Context)
	SwitchState(*gin.Context)
	Send2(*gin.Context)
}

type AgentImpl struct {
	enable bool
	trx    client.TrxClient
	srv    server.Server
	l      *utils.Logger
}

func NewAgentImpl(log mylog.Logging, cli client.TrxClient, s server.Server) Agent {
	a := AgentImpl{
		enable: true,
		trx:    cli,
		srv:    s,
		l:      log.GetLogger("agent"),
	}

	a.l.Info("Init...")

	return &a
}

func (a *AgentImpl) SwitchOn(c *gin.Context) {
	a.enable = true
	c.JSON(200, responseSuccess("ok"))
}

func (a *AgentImpl) SwitchOff(c *gin.Context) {
	a.enable = false
	c.JSON(200, responseSuccess("ok"))
}

func (a *AgentImpl) SwitchState(c *gin.Context) {
	msg := "On"
	if !a.enable {
		msg = "Off"
	}
	c.JSON(200, responseSuccess(msg))
}

type TxSendReq struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Token  string `json:"token"`
	Amount int64  `json:"amount"`
	Note   string `json:"note"`
}

func getSwitch() bool {
	url := "http://47.242.130.201:9998/switch/get"

	req, _ := http.NewRequest("GET", url, nil)
	//req.Close = true //短链接，防止重用tcp连接。期望解决read: connection reset by peer
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	jsonValue, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	var ret SwtichT

	err = json.Unmarshal(jsonValue, &ret)
	if err != nil {
		return false
	}

	if !ret.Switch {
		return false
	}

	return true
}

type SwtichT struct {
	Switch bool `json:"switch"`
}

func (a *AgentImpl) Send2(c *gin.Context) {
	ip := c.ClientIP()
	if ip != "127.0.0.1" {
		a.l.Error("Stop! Ip was banned")
		c.JSON(200, ErrClientIP)
		return
	}

	sw := getSwitch()
	if !sw {
		a.l.Error("Stop! Switch is close")
		c.JSON(200, ErrSwitch)
		return
	}

	var req TxSendReq
	bodyBytes, _ := io.ReadAll(c.Request.Body)

	err := json.Unmarshal(bodyBytes, &req)
	if err != nil {
		a.l.Error(err.Error())
		c.JSON(200, ErrIncorrectFormat)
		return
	}

	hash, err := a.trx.Transfer(req.From, req.To, req.Amount, req.Token)
	if err != nil {
		a.l.Error(err.Error())
		c.JSON(200, ErrInternalError)
		return
	}
	c.JSON(200, responseSuccess(hash))
}

func (a *AgentImpl) RegisterRouters() {
	a.l.Info("resister handler...")
	a.srv.RegisterHandler("get", "/switch/on", a.SwitchOn)
	a.srv.RegisterHandler("get", "/switch/off", a.SwitchOff)
	a.srv.RegisterHandler("get", "/switch/state", a.SwitchState)
	a.srv.RegisterHandler("post", "/send2", a.Send2)
}

var AgentModule = fx.Options(fx.Provide(NewAgentImpl))
var RegisterRouters = fx.Options(fx.Invoke(Agent.RegisterRouters))
