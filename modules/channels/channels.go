package channels

import (
	"github.com/dafsic/gambler/lib/mylog"
	"github.com/dafsic/gambler/utils"
	"go.uber.org/fx"
)

type ChanManager interface {
	GetChan(string) chan interface{}
	//CreateChan(string, int) chan interface{}
	CloseChan(string)
}

type ChanManagerImpl struct {
	chans map[string]chan interface{}
	//rwlock sync.RWMutex
	l *utils.Logger
}

func (cm *ChanManagerImpl) GetChan(name string) chan interface{} {
	return cm.chans[name]
}

// func (cm *ChanManagerImpl) CreateChan(name string, size int) chan interface{} {
// 	cm.chans[name] = make(chan interface{}, size)
// 	return cm.chans[name]
// }

func (cm *ChanManagerImpl) CloseChan(name string) {
	c := cm.chans[name]
	close(c)
	//delete(cm.chans, name)
}

func NewChanManager(log mylog.Logging) ChanManager {
	cs := make(map[string]chan interface{}, 8)
	c := &ChanManagerImpl{
		chans: cs,
		l:     log.GetLogger("channel"),
	}

	c.l.Info("Init...")
	// 下注交易的hash
	c.chans["bethash"] = make(chan interface{}, 8)
	// 需要下注的数量
	c.chans["bet"] = make(chan interface{}, 8)
	// 打包进区块的所有交易hash、区块高度、区块hash
	c.chans["block"] = make(chan interface{}, 128)

	return c
}

var ChanManagerModule = fx.Options(fx.Provide(NewChanManager))
