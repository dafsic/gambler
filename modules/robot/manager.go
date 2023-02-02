package robot

import (
	"context"
	"encoding/json"

	"github.com/bwmarrin/snowflake"
	"github.com/dafsic/gambler/lib/mylog"
	"github.com/dafsic/gambler/modules"
	"github.com/dafsic/gambler/modules/channels"
	"github.com/dafsic/gambler/modules/client"
	"github.com/dafsic/gambler/modules/server"
	"github.com/dafsic/gambler/modules/store"
	"github.com/dafsic/gambler/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type RobotManager interface {
	RegisterRouters()
	Working()
}

type RobotManagerImpl struct {
	// 管理所有的robot，不是频繁创建和释放的对象，不需要池
	robots map[string]Robot
	db     store.Store
	l      *utils.Logger
	trx    client.TrxClient
	snow   *snowflake.Node
	srv    server.Server
	blockC chan interface{}
	qc     chan bool
}

func NewRobotManager(lc fx.Lifecycle, log mylog.Logging, chanMgr channels.ChanManager, cli client.TrxClient, db store.Store, s server.Server) (RobotManager, error) {
	m := &RobotManagerImpl{
		l:      log.GetLogger("manager"),
		blockC: chanMgr.GetChan("block"),
		trx:    cli,
		robots: make(map[string]Robot, 32),
		db:     db,
		srv:    s,
		qc:     make(chan bool, 1),
	}

	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, err
	}
	m.snow = node

	lc.Append(fx.Hook{
		// app.start调用
		OnStart: func(ctx context.Context) error {
			// 这里不能阻塞
			go m.Working()
			return nil
		},
		// app.stop调用，收到中断信号的时候调用app.stop
		OnStop: func(ctx context.Context) error {
			go m.Stop()
			return nil
		},
	})

	return m, nil
}

func (m *RobotManagerImpl) Working() {
	for {
		select {
		case <-m.qc:
			m.l.Info("exit...")
			return
		default:
			b := <-m.blockC
			block := b.(*modules.Block)
			for _, v := range m.robots {
				v.ReceiveBlock(block)
			}
		}
	}
}

func (m *RobotManagerImpl) Stop() {
	m.qc <- true
}

type RobotConfig struct {
	Pool      string `json:"pool"`
	Token     string `json:"token"`
	StartNum  int    `json:"start_num"`
	NumOfBets int    `json:"num_of_bets"`
	OddChips  []int  `json:"odd_chips"`
	EvenChips []int  `json:"even_chips"`
	TP        int    `json:"take_profit"`
	SL        int    `json:"stop_loss"`
}

func (m *RobotManagerImpl) CreateRobot(c *gin.Context) {
	var req RobotConfig
	body, _ := c.Get("body")

	err := json.Unmarshal(body.([]byte), &req)
	if err != nil {
		m.l.Error(err.Error())
		c.JSON(200, ErrIncorrectFormat)
		return
	}

	poolInfo, err := m.db.GetPoolByAddr(req.Pool, req.Token)
	if err != nil {
		m.l.Error(err.Error())
		c.JSON(200, ErrInternalError)
		return
	}

	rid := m.snow.Generate().String()

	addr, key, err := m.trx.Generateaddress()
	if err != nil {
		m.l.Error(err.Error())
		c.JSON(200, ErrInternalError)
		return
	}

	para := store.RobotParameter{
		Rid:       rid,
		PoolId:    poolInfo.Id,
		Addr:      addr,
		Key:       key,
		StartNum:  req.StartNum,
		NumOfBets: req.NumOfBets,
		TP:        req.TP,
		SL:        req.SL,
		EvenChips: req.EvenChips,
		OddChips:  req.EvenChips,
	}

	err = m.db.CreateRobotParameter(&para)
	if err != nil {
		m.l.Error(err.Error())
		c.JSON(200, ErrInternalError)
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "success", "id": rid, "addr": addr, "version": ver})
	return
}

func (m *RobotManagerImpl) RunRobot(c *gin.Context) {
	rid, exist := c.GetQuery("id")
	if !exist {
		m.l.Error("missing parameter")
		c.JSON(200, ErrParameter)
		return
	}

	para, err := m.db.GetRobotParameter(rid)
	if err != nil { //ErrNotFound 同样处理
		m.l.Error(err.Error())
		c.JSON(200, ErrInternalError)
		return
	}

	poolInfo, err := m.db.GetPoolById(para.PoolId)
	if err != nil {
		m.l.Error(err.Error())
		c.JSON(200, ErrInternalError)
		return
	}

	r := NewRobot(para, poolInfo, m.l, m.trx)
	m.robots[rid] = r
	r.Working()

	para.State = 1
	err = m.db.UpdateRobotParameter(para.Rid, para)
	if err != nil {
		m.l.Error(err.Error())
		c.JSON(200, ErrInternalError)
		return
	}

	c.JSON(200, responseSuccess("success"))
	return
}

func (m *RobotManagerImpl) StopRobot(c *gin.Context) {
	rid, exist := c.GetQuery("id")
	if !exist {
		m.l.Error("missing parameter")
		c.JSON(200, ErrParameter)
		return
	}

	r, ok := m.robots[rid]
	if !ok {
		m.l.Errorf("robot[%s] not found\n", rid)
		c.JSON(200, ErrNotFound)
		return
	}

	m.robots[rid].Exit()  //退出线程
	delete(m.robots, rid) //释放内存

	r.PrintParameter().State = 0
	err := m.db.UpdateRobotParameter(rid, r.PrintParameter())
	if err != nil {
		m.l.Error(err.Error())
		c.JSON(200, ErrInternalError)
		return
	}

	c.JSON(200, responseSuccess("success"))
	return
}

func (m *RobotManagerImpl) UpdateRobotParameter(c *gin.Context) {
	rid, exist := c.GetQuery("id")
	if !exist {
		m.l.Error("missing parameter")
		c.JSON(200, ErrParameter)
		return
	}

	var req RobotConfig
	body, _ := c.Get("body")

	err := json.Unmarshal(body.([]byte), &req)
	if err != nil {
		m.l.Error(err.Error())
		c.JSON(200, ErrIncorrectFormat)
		return
	}

	poolInfo, err := m.db.GetPoolByAddr(req.Pool, req.Token)
	if err != nil {
		m.l.Error(err.Error())
		c.JSON(200, ErrInternalError)
		return
	}

	para := store.RobotParameter{
		Rid:       rid,
		PoolId:    poolInfo.Id,
		StartNum:  req.StartNum,
		NumOfBets: req.NumOfBets,
		TP:        req.TP,
		SL:        req.SL,
		EvenChips: req.EvenChips,
		OddChips:  req.EvenChips,
	}

	err = m.db.UpdateRobotParameter(rid, &para)
	if err != nil {
		m.l.Error(err.Error())
		c.JSON(200, ErrInternalError)
		return
	}

	c.JSON(200, responseSuccess("success"))
	return
}

// 由于机器人的池子是可变的，所以查询余额，不能预设就是当前配置指定token的余额
func (m *RobotManagerImpl) GetBalance(c *gin.Context) {
	addr, e1 := c.GetQuery("addr")
	token, e2 := c.GetQuery("token")
	if !e1 || !e2 {
		m.l.Error("missing parameter")
		c.JSON(200, ErrParameter)
		return
	}

	amount, err := m.trx.GetBalance(addr, token)
	if err != nil {
		m.l.Error(err.Error())
		c.JSON(200, ErrInternalError)
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "success", "balance": float64(amount) / 1000000, "version": ver})
	return
}

func (m *RobotManagerImpl) RegisterRouters() {
	m.l.Info("resister handler...")
	m.srv.RegisterHandler("post", "/robot/create", m.CreateRobot)
	m.srv.RegisterHandler("get", "/robot/run", m.RunRobot)
	m.srv.RegisterHandler("get", "/robot/stop", m.StopRobot)
	m.srv.RegisterHandler("post", "/robot/update", m.UpdateRobotParameter)
	m.srv.RegisterHandler("get", "/robot/balance", m.GetBalance)
}

var SchedModule = fx.Options(fx.Provide(NewRobotManager))
var RegisterRouters = fx.Options(fx.Invoke(RobotManager.RegisterRouters))
