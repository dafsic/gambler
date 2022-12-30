package listent

import (
	"context"
	"encoding/json"

	"github.com/dafsic/gambler/config"
	"github.com/dafsic/gambler/lib/mylog"
	"github.com/dafsic/gambler/modules"
	"github.com/dafsic/gambler/modules/channels"
	"github.com/dafsic/gambler/utils"
	kafka "github.com/segmentio/kafka-go"
	"go.uber.org/fx"
)

type Listener interface {
	Working()
	Stop()
}

type ListenerImpl struct {
	kafka_brokers []string
	kafka_topic   string
	wc            chan interface{}
	qc            chan bool
	l             *utils.Logger
}

func NewListener(lc fx.Lifecycle, log mylog.Logging, cfg config.ConfigI, chanMgr channels.ChanManager) Listener {

	listener := &ListenerImpl{
		l:             log.GetLogger("listener"),
		wc:            chanMgr.GetChan("block"),
		qc:            make(chan bool, 1),
		kafka_brokers: cfg.GetElem("kafkaBrokers").([]string),
		kafka_topic:   cfg.GetElem("kafkaTopic").(string),
	}
	listener.l.Info("Init...")
	lc.Append(fx.Hook{
		// app.start调用
		OnStart: func(ctx context.Context) error {
			// 这里不能阻塞
			go listener.Working()
			return nil
		},
		// app.stop调用，收到中断信号的时候调用app.stop
		OnStop: func(ctx context.Context) error {
			go listener.Stop()
			return nil
		},
	})

	return listener
}

func (listerner *ListenerImpl) Working() {
	// make a new reader that consumes from topic-A, partition 0, at last offset
	// 不写GroupID，则不读历史数据
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   listerner.kafka_brokers,
		Topic:     listerner.kafka_topic,
		Partition: 0,
		MinBytes:  256,  // 1KB
		MaxBytes:  10e6, // 10MB
	})
	r.SetOffset(kafka.LastOffset)

	for {
		select {
		case <-listerner.qc:
			listerner.l.Info("Quit...")
			if err := r.Close(); err != nil {
				listerner.l.Error(err.Error())
				return
			}
		default:
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				listerner.l.Error(err.Error())
				continue
			}
			// 解析message，输出交易信息
			var block modules.Block
			err = json.Unmarshal(m.Value, &block)
			if err != nil {
				listerner.l.Error(err.Error())
				continue
			}
			listerner.l.Infof("Receive a block --> ts:%d, hash: %s, height: %d,txs:%d\n", block.Ts, block.BlockHash, block.Height, len(block.Txs))
			listerner.wc <- &block
		}
	}
}

func (listener *ListenerImpl) Stop() {
	listener.qc <- true
}

var ListenModule = fx.Options(fx.Invoke(NewListener))
