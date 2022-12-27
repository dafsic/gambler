package sched

import (
	"github.com/dafsic/gambler/lib/mylog"
	"github.com/dafsic/gambler/modules"
	"github.com/dafsic/gambler/modules/channels"
	"go.uber.org/fx"
)

type Scheduler interface {
}

var evenBetAmount = []int{23, 61, 143, 311, 657, 1367, 2823, 5813, 11949}
var oddBetAmount = []int{22, 60, 140, 306, 646, 1344, 2776, 5716, 11752}

type SchedulerImpl struct {
	lastBlockHash   string
	lastBlockHeight int64
	lastBetHash     string
	blockC          chan interface{}
	betC            chan interface{}
	betHashC        chan interface{}
	betCounter      [2]int64 //0:偶数下注次数，1:奇数下注次数
	blockCounter    [2]int   //0:区块hash是偶数的连续次数，1:连续奇数的次数
	status          string   //"odd":下了偶数的注, "even":下了奇数的注, "":没下注
}

func NewSchedulerImpl(lc fx.Lifecycle, log mylog.Logging, chanMgr channels.ChanManager) Scheduler {
	s := &SchedulerImpl{
		blockC:   chanMgr.GetChan("block"),
		betC:     chanMgr.GetChan("bet"),
		betHashC: chanMgr.GetChan("bethash"),
	}
	return s
}

func (s *SchedulerImpl) Working() {
	for {
		select {
		//优先级队列
		case betHash := <-s.betHashC:
			s.lastBetHash = betHash.(string)
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

	if s.blockCounter[0] > 7 && s.status == "" {
		s.betC <- evenBetAmount[s.betCounter[1]]
		s.status = "even"
		s.betCounter[1] += 1
	}

	if s.blockCounter[1] > 7 && s.status == "" {
		s.betC <- oddBetAmount[s.betCounter[0]]
		s.status = "odd"
		s.betCounter[0] += 1
	}

	// 如果下单失败，拿不到下注hash，就直接返回，观察到一直不下注，就会看到日志中的下注失败的错误
	if s.status == "" || s.lastBetHash == "" {
		return
	}

	for _, tx := range block.Txs {
		if tx == s.lastBetHash { //下注交易包含在区块中
			if (s.status == "odd" && isOdd) || (s.status == "even" && !isOdd) { //中了
				s.status = ""
				s.lastBetHash = ""
				s.betCounter[0] = 0
				s.betCounter[1] = 0
			} else { //没中
				if s.blockCounter[0] > 7 { //没中但偶数连续情况并未被打断，不管是否跳块，继续下注奇数
					s.betC <- evenBetAmount[s.betCounter[1]]
					s.betCounter[1] += 1
				} else if s.betCounter[1] > 7 { //没中但奇数连续情况并未被打断，不管是否跳块，继续下注偶数
					s.betC <- oddBetAmount[s.betCounter[0]]
					s.betCounter[0] += 1
				}
				//没中，但连续奇偶都不足7，说明在跳块中中了，那就认赔了
			}

			return
		}
	}

	// 交易没包含在此区块中，则什么都不做，只统计连续奇偶情况
}

var oddMap = map[byte]bool{'0': true, '1': false, '2': true, '3': false, '4': true, '5': false, '6': true, '7': false, '8': true, '9': false}

// IsOddNum 判断一个字符串最后一个数字是偶数吗
func IsOddNum(h string) bool {
	l := len(h)
	for i := 1; i <= l; i++ {
		if h[l-i] > '9' {
			continue
		}
		return oddMap[h[l-i]]
	}
	return true
}
