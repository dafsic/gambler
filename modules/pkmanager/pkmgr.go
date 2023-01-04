package pkmanager

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/dafsic/gambler/config"
	"go.uber.org/fx"
)

type PKManager interface {
	GetPrivateKey(addr string) (string, bool)
}

type PKManagerImpl struct {
	Keys map[string]string
}

func NewPKManager(cfg config.ConfigI) (PKManager, error) {
	pkm := PKManagerImpl{
		Keys: make(map[string]string, 8),
	}
	f := cfg.GetElem("keys").(string)
	err := pkm.parseKeysFile(f)
	if err != nil {
		return nil, err
	}

	return &pkm, nil
}

func (pkm *PKManagerImpl) GetPrivateKey(addr string) (string, bool) {
	key, ok := pkm.Keys[addr]
	return key, ok
}

func (pkm *PKManagerImpl) parseKeysFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		s := strings.TrimSpace(string(line))
		if len(s) == 0 || strings.HasPrefix(s, "#") {
			continue
		}
		index := strings.Index(s, ":")
		if index < 0 {
			continue
		}
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		pkm.Keys[key] = value
	}
	return nil
}

var PKManagerModule = fx.Options(fx.Provide(NewPKManager))
