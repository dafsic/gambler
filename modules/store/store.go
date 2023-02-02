// robot 参数读写,复杂后可以引入gorm
package store

type Store interface {
	GetRobotParameter(rid string) (*RobotParameter, error)
	CreateRobotParameter(para *RobotParameter) error
	// UpdateRobotParameter 不能更改addr和key
	UpdateRobotParameter(rid string, para *RobotParameter) error
	GetPoolById(id int) (*PoolInfo, error)
	GetPoolByAddr(addr, token string) (*PoolInfo, error)
}

type RobotParameter struct {
	Rid       string
	PoolId    int    // pool id
	Addr      string // robot 地址
	Key       string // private key
	StartNum  int    // 第n次连续开始下注
	NumOfBets int    // 最多连续下n次
	TP        int    // 止盈
	SL        int    // 止损
	OddChips  []int  // 奇数筹码
	EvenChips []int  // 偶数筹码
	State     int    // 状态
}

func GetRobotParameter(rid string) (*RobotParameter, error) {
	var r RobotParameter
	return &r, nil
}

func CreateParameter(r *RobotParameter) error {
	return nil
}

func UpdateRobotParameter(rid string, para *RobotParameter) error {
	return nil
}

type PoolInfo struct {
	Id        int
	Addr      string
	Refund    string
	Token     string
	MinAmount int
	MaxAmount int
	State     int
}

func GetPoolInfo(int) (*PoolInfo, error) {
	var p PoolInfo
	return &p, nil
}
