// 日志的level等信息应该从配置文件里读取，配置文件的路径应该从环境变量或者命令行参数里获取
package mylog

import (
	"io"
	"os"
	"sync"

	"github.com/dafsic/gambler/config"
	"github.com/dafsic/gambler/utils"
	"go.uber.org/fx"
)

type Logging interface {
	GetLogger(name string) *utils.Logger
}

type LoggingImpl struct {
	mux     sync.Mutex
	lvl     string
	output  io.Writer
	loggers map[string]*utils.Logger
}

func (l *LoggingImpl) GetLogger(name string) *utils.Logger {
	l.mux.Lock()
	defer l.mux.Unlock()
	i, ok := l.loggers[name]
	if !ok {
		i = utils.NewLogger(l.output, name, utils.LogLevelFromString(l.lvl), utils.Ldefault)
		l.loggers[name] = i
	}
	return i
}

// var once sync.Once
// var l Logging

// func NewMylog(cfg config.ConfigI) Logging {
// 	once.Do(func() {
// 		var t LoggingImpl
// 		t.output = os.Stdout
// 		t.lvl = cfg.GetElem("logLevel").(string)
// 		t.loggers = make(map[string]*utils.Logger, 8)
// 		l = &t
// 	})
// 	return l
// }

func NewMylog(cfg config.ConfigI) Logging {
	t := &LoggingImpl{
		output:  os.Stdout,
		lvl:     cfg.GetElem("logLevel").(string),
		loggers: make(map[string]*utils.Logger, 8),
	}

	return t
}

var Module = fx.Options(fx.Provide(NewMylog))
