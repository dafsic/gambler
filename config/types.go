package config

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"reflect"
	"strings"
)

type ConfigI interface {
	GetElem(e string) interface{}
}

// robot config
type GamblerCfg struct {
	LogLevel     string   `toml:"loglevel"`
	TrxNode      string   `toml:"trxnode"` //trx节点
	Addr         string   `toml:"addr"`    //下注地址
	Pool         string   `toml:"pool"`    //池
	Refund       string   `toml:"refund"`  //回款地址
	Keys         string   `toml:"keys"`    //存放地址->私钥的txt文件
	Listen       string   `toml:"listen"`  //http监听地址(transfer使用)
	Token        string   `toml:"token"`   //token种类(trx,usdt...)
	KafkaTopic   string   `toml:"kafkaTopic"`
	KafkaBrokers []string `toml:"kafkaBrokers"`
}

func (a *GamblerCfg) GetElem(e string) interface{} {
	var cfg interface{}
	rt := reflect.TypeOf(*a)
	rv := reflect.ValueOf(*a)

	fieldNum := rt.NumField()
	for i := 0; i < fieldNum; i++ {
		if strings.ToUpper(rt.Field(i).Name) == strings.ToUpper(e) {
			cfg = rv.FieldByName(rt.Field(i).Name).Interface()
			break
		}
	}
	return cfg
}

func NewRobotCfg(ctx *cli.Context) (ConfigI, error) {
	cp := ctx.String("config")
	c, err := FromFile(cp, DefaultGamblerCfg())
	if err != nil {
		return nil, err
	}
	cfg, ok := c.(*GamblerCfg)
	if !ok {
		return nil, fmt.Errorf("invalid config for assistant, got: %T", c)
	}
	//fmt.Printf("config---->topic:%s,Brokers:%v\n", cfg.KafkaTopic, cfg.KafkaBrokers)
	return cfg, nil
}

var CfgModule = fx.Options(fx.Provide(NewRobotCfg))
