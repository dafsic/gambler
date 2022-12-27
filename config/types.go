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
	LogLevel string `toml:"loglevel"`
	TrxNode  string `toml:"trxnode"`
	Addr     string `toml:"addr"`
	Pk       string `toml:"pk"`
	Pool     string `toml:"pool"`
	Refund   string `toml:"refund"`
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
	pk := ctx.String("pk")
	if pk != "" {

	}

	return cfg, nil
}

var CfgModule = fx.Options(fx.Provide(NewRobotCfg))
