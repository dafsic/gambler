package sched

import (
	"context"
	"fmt"

	"github.com/dafsic/gambler/config"
	"github.com/dafsic/gambler/lib/mylog"
	"github.com/dafsic/gambler/modules"
	"github.com/dafsic/gambler/modules/channels"
	"github.com/dafsic/gambler/modules/client"
	"github.com/dafsic/gambler/utils"
	"go.uber.org/fx"
)

type Scheduler interface {
	Working()
	Stop()
}

var trxBase int64 = 1000000
var usdtBase int64 = 1000000
var betTriggerNum = 3

// 30000
var minBalance = 1

// 0:下偶数注，1:下奇数注
//var BetAmount = [2][]int64{{20, 40, 80, 160, 320, 640, 1280, 2560, 5120}, {21, 41, 81, 161, 321, 641, 1281, 2561, 5121}}

// usdt用
var BetAmount = [2][]int64{{10, 20, 40, 80, 160, 320}, {11, 21, 41, 81, 161, 321}}

type SchedulerImpl struct {
	lastBetHash  string
	refund       string //回款地址
	pool         string
	addr         string //下注的地址
	token        string //币种
	blockC       chan interface{}
	betC         chan interface{}
	betHashC     chan interface{}
	qc           chan bool
	betCounter   [2]int64 //0:偶数下注次数, 1:奇数下注次数
	blockCounter [2]int   //0:区块hash是偶数的连续次数, 1:连续奇数的次数
	status       int      //0:下了偶数的注, 1:下了奇数的注, 其他:无效
	trx          client.TrxClient
	l            *utils.Logger
	isRefund     bool
}

func NewScheduler(lc fx.Lifecycle, log mylog.Logging, cfg config.ConfigI, chanMgr channels.ChanManager, cli client.TrxClient) Scheduler {
	s := &SchedulerImpl{
		pool:     cfg.GetElem("pool").(string),
		addr:     cfg.GetElem("addr").(string),
		refund:   cfg.GetElem("refund").(string),
		token:    cfg.GetElem("token").(string),
		blockC:   chanMgr.GetChan("block"),
		betC:     chanMgr.GetChan("bet"),
		betHashC: chanMgr.GetChan("bethash"),
		qc:       make(chan bool, 1),
		trx:      cli,
		l:        log.GetLogger("sched"),
		isRefund: true,
		status:   2,
	}
	s.l.Info("Init...")
	lc.Append(fx.Hook{
		// app.start调用
		OnStart: func(ctx context.Context) error {
			// 这里不能阻塞
			go s.Working()
			return nil
		},
		// app.stop调用，收到中断信号的时候调用app.stop
		OnStop: func(ctx context.Context) error {
			go s.Stop()
			return nil
		},
	})

	return s
}

func (s *SchedulerImpl) Stop() {
	s.qc <- true
}

func (s *SchedulerImpl) Working() {
	for {
		select {
		//优先级队列
		case <-s.qc:
			return
		default:
			block := <-s.blockC
			s.dealBlock(block.(*modules.Block))
		}
	}
}

func (s *SchedulerImpl) dealBlock(block *modules.Block) {
	isOdd := IsOddNum(block.BlockHash)
	if isOdd {
		s.blockCounter[0] += 1
		s.blockCounter[1] = 0
	} else {
		s.blockCounter[0] = 0
		s.blockCounter[1] += 1
	}

	s.l.Infof("odd:%d,even:%d\n", s.blockCounter[0], s.blockCounter[1])

	//如果之前中了，要等回款后才继续下注
	if !s.isRefund {
		r, e := IsRefund(s.refund, s.addr, block.Ts-30000, block.Ts, s.token)
		if e != nil {
			s.l.Error(e.Error())
		}
		s.isRefund = r
		s.l.Infof("Refund:%t\n", r)
		return
	}

	//尝试开启第一次下注
	success, err := s.tryBet(true)
	if err != nil {
		s.l.Error(err.Error())
		return
	}

	if success {
		return
	}

	//没下注就不需要处理块
	if s.lastBetHash == "" {
		return
	}

	for _, tx := range block.Txs {
		if tx == s.lastBetHash { //下注交易包含在区块中
			if (s.status == 0 && isOdd) || (s.status == 1 && !isOdd) { //中了
				s.l.Infof("中了,block:%s\n", block.BlockHash)
				s.status = 2
				s.lastBetHash = ""
				s.betCounter[0] = 0
				s.betCounter[1] = 0
				s.isRefund = false
			} else {
				s.l.Infof("没中,block:%s\n", block.BlockHash)
				//没中，但奇偶连续情况并未打断的情况下，不用管是否跳块，继续下注
				success, err = s.tryBet(false)
				if err != nil {
					s.l.Error(err.Error())
					return
				}

				if success {
					//如果跳块超过7个，并且包含连续超过7个奇偶的时候，就会出现不是从起始下注金额下注的情况
					//已在tryBet中处理
					s.l.Info("没中,继续下了一注...")
				} else {
					s.l.Info("没中,但在跳块中中了,或者已连续投注9次了")
					s.lastBetHash = ""
					s.betCounter[0] = 0
					s.betCounter[1] = 0
					s.status = 2
				}
			}
			return
		}
	}

	// 交易没包含在此区块中，则什么都不做，只统计连续奇偶情况
	s.l.Info("跳块")
}

func (s *SchedulerImpl) tryBet(first bool) (bool, error) {
	//已经下过第一次注了，不用下了
	if first && s.lastBetHash != "" {
		return false, nil
	}

	var base int64
	switch s.token {
	case "trx":
		base = trxBase
	case "usdt":
		base = usdtBase
	default:
		return false, nil
	}

	maxBetNum := len(BetAmount[0])
	for i, v := range s.blockCounter {
		turn := (i + 1) % 2
		//s.l.Infof("下注条件,i=%d,v=%d,betCounter[%d]=%d\n", i, v, turn, s.betCounter[turn])
		if v > betTriggerNum && s.betCounter[turn] < int64(maxBetNum) {
			//连续跳块超过7次的时候，可能会出现这种反转下注情况
			if i == s.status {
				return false, nil
			}
			balance, err := s.trx.GetBalance(s.addr)
			if err != nil {
				return false, fmt.Errorf("%w%s", err, utils.LineNo())
			}

			if balance < int64(minBalance)*base {
				return false, fmt.Errorf("insufficient balance < %d,%s", balance, utils.LineNo())
			}

			hash, err := s.trx.Transfer(s.addr, s.pool, BetAmount[turn][s.betCounter[turn]]*base, s.token)
			if err != nil {
				return false, fmt.Errorf("%w%s", err, utils.LineNo())
			}

			s.l.Infof("下注:%d,hash:%s\n", BetAmount[turn][s.betCounter[turn]], hash)
			s.lastBetHash = hash
			s.betCounter[turn] += 1
			s.status = turn
			return true, nil
		}
	}

	return false, nil
}

// var SchedModule = fx.Options(fx.Provide(NewScheduler))
var SchedModule = fx.Options(fx.Invoke(NewScheduler))
