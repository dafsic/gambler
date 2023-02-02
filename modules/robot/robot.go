package robot

import (
	"context"
	"fmt"
	"strings"

	"github.com/dafsic/gambler/modules"
	"github.com/dafsic/gambler/modules/client"
	"github.com/dafsic/gambler/modules/store"
	"github.com/dafsic/gambler/utils"
)

var trxBase int64 = 1000000
var usdtBase int64 = 1000000

// 无法限制直接赋
type sta int

func (s *sta) Set(a sta) {
	if a != EVEN && a != ODD && a != INVALID {
		panic("sta error")
	}
	*s = sta(a)
}

func (s *sta) Turn() sta {
	switch *s {
	case EVEN:
		return ODD
	case ODD:
		return EVEN
	default:
		return INVALID
	}
}

const (
	EVEN sta = iota
	ODD
	INVALID
)

type Robot interface {
	Working()                              // 运行
	Exit()                                 // 停止
	PrintParameter() *store.RobotParameter // 使用的参数
	ReceiveBlock(*modules.Block)           // 接收最新区块
}

type RobotImpl struct {
	para         *store.RobotParameter
	pool         *store.PoolInfo
	cancel       context.CancelFunc
	lastBetHash  string
	blockC       chan *modules.Block
	betCounter   int         //连续下注次数
	blockCounter map[sta]int // EVEN:偶数区块连续次数, ODD:奇数区块连续次数
	state        sta         // EVEN:下了偶数的注, ODD:下了奇数的注, INVALID:无效
	trx          client.TrxClient
	l            *utils.Logger
	isRefund     bool
}

func NewRobot(rp *store.RobotParameter, pi *store.PoolInfo, log *utils.Logger, cli client.TrxClient) Robot {
	r := &RobotImpl{
		para:         rp,
		pool:         pi,
		blockCounter: make(map[sta]int, 2),
		//betCounter:   make(map[sta]int, 2),
		blockC:   make(chan *modules.Block, 1), //只有1个缓冲，阻塞了就停止机器人
		trx:      cli,
		l:        log,
		isRefund: true,
	}
	r.state.Set(INVALID)

	r.l.Infof("Robot[%s] Working...\n", rp.Rid)
	return r
}

// 必须立刻处理，不能等待
func (r *RobotImpl) ReceiveBlock(b *modules.Block) {
	select {
	case r.blockC <- b:
	default:
		r.l.Errorf("robot[%s] channel is full,stop!", r.para.Rid)
		r.cancel()
	}
}

func (r *RobotImpl) Working() {
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	for {
		select {
		//优先级队列
		case <-ctx.Done():
			close(r.blockC)
			return
		default:
			block := <-r.blockC
			r.dealBlock(block)
		}
	}
}

func (r *RobotImpl) Exit() {
	r.cancel()
}

func (r *RobotImpl) PrintParameter() *store.RobotParameter {
	return r.para
}

func (r *RobotImpl) dealBlock(block *modules.Block) {
	evenBlock := IsEvenNum(block.BlockHash)
	if evenBlock {
		r.blockCounter[EVEN] += 1
		r.blockCounter[ODD] = 0
	} else {
		r.blockCounter[EVEN] = 0
		r.blockCounter[ODD] += 1
	}

	r.l.Infof("robot[%s] odd:%d,even:%d\n", r.para.Rid, r.blockCounter[ODD], r.blockCounter[EVEN])

	//如果之前中了，要等回款后才继续下注
	if !r.isRefund {
		b, e := IsRefund(r.pool.Refund, r.para.Addr, block.Ts-60000, block.Ts, r.pool.Token)
		if e != nil {
			r.l.Error(e.Error())
		}
		r.isRefund = b
		r.l.Infof("robot[%s] Refund:%t\n", r.para.Rid, b)
		return
	}

	//尝试开启第一次下注(新一轮)
	success, err := r.tryBet(true)
	if err != nil {
		r.l.Errorf("robot[%s],err:%s\n", r.para.Rid, err.Error())
		return
	}

	// 第一次下注成功or没下过注，就不需要处理区块
	if success || r.state == INVALID {
		return
	}

	for _, tx := range block.Txs {
		if tx == r.lastBetHash { //下注交易包含在区块中
			if (r.state == ODD && evenBlock) || (r.state == EVEN && !evenBlock) { //中了
				r.l.Infof("robot[%s] 中了,block:%s\n", r.para.Rid, block.BlockHash)
				r.state.Set(INVALID)
				r.lastBetHash = ""
				r.betCounter = 0
				r.isRefund = false
			} else {
				r.l.Infof("robot[%s] 没中,block:%s\n", block.BlockHash)
				//没中，但奇偶连续情况并未打断的情况下，不用管是否跳块，继续下注
				success, err = r.tryBet(false)
				if err != nil {
					// 下注出错，这一轮还可以接着下
					r.l.Errorf("robot[%s] err:%s\n", r.para.Rid, err.Error())
					return
				}
				if !success {
					//没错,但不满足下注条件了,需要重置
					r.l.Infof("robot[%s] 没中,但在跳块中中了,或者超过最大下注次数,或者止盈止损了\n", r.para.Rid)
					r.lastBetHash = ""
					r.betCounter = 0
					r.state.Set(INVALID)
				}
			}
			return
		}
	}

	// 交易没包含在此区块中，则什么都不做，只统计连续奇偶情况
	r.l.Infof("robot[%s] 跳块\n", r.para.Rid)
}

// tryBet 如果到达下注条件就尝试去下注
func (r *RobotImpl) tryBet(first bool) (bool, error) {
	// 已经下过第一注了，不用再下第一注了
	if first && r.state != INVALID {
		return false, nil
	}

	var base int64
	switch strings.ToUpper(r.pool.Token) {
	case "TRX":
		base = trxBase
	case "USDT":
		base = usdtBase
	default:
		return false, fmt.Errorf("不支持的token:%s,%s", r.pool.Token, utils.LineNo())
	}

	for s, c := range r.blockCounter {
		turn := s.Turn()
		//r.l.Infof("下注条件,s=%d,c=%d,betCounter[%d]=%d\n", s, c, turn, r.betCounter[turn])
		if c >= r.para.StartNum && r.betCounter < r.para.NumOfBets {
			//连续跳块超过7次的时候，可能会出现这种反转下注情况
			if s == r.state {
				r.l.Warnf("robot[%s] 下注失败,不允许反转下注\n", r.para.Rid)
				return false, nil
			}
			balance, err := r.trx.GetBalance(r.para.Addr, r.pool.Token)
			if err != nil {
				return false, fmt.Errorf("%w%s", err, utils.LineNo())
			}

			var amount int64
			switch s {
			case EVEN:
				amount = int64(r.para.OddChips[r.betCounter+1])
			case ODD:
				amount = int64(r.para.EvenChips[r.betCounter+1])
			}

			if balance-amount*base <= int64(r.para.SL)*base || balance >= int64(r.para.TP) {
				r.l.Warnf("robot[%s] 下注失败,止盈或者止损,banlance:%d\n", r.para.Rid, balance/base)
				return false, nil
			}

			hash, err := r.trx.Transfer(r.para.Addr, r.pool.Addr, r.para.Key, amount*base, r.pool.Token)
			if err != nil {
				return false, fmt.Errorf("%w%s", err, utils.LineNo())
			}

			r.l.Infof("robot[%s] 下注成功:%d,hash:%s\n", amount, hash)
			r.lastBetHash = hash
			r.betCounter += 1
			r.state.Set(turn)
			return true, nil
		}
	}

	//到这里,说明奇偶连续次数不够,如果不是第一次
	//就代表在跳块中中了
	return false, nil
}
