// 管理robot参数模块，依赖存储模块
package robot

type parameter struct {
	addr   string // robot 地址
	key    string // private key
	token  string // trx/usdr
	pool   string // pool 地址
	refund string // 回款地址
}

func GetParameter(rid string) *parameter {

}
