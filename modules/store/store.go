// robot 参数读写,复杂后可以引入gorm
package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dafsic/gambler/config"
	"github.com/dafsic/gambler/lib/mylog"
	"github.com/dafsic/gambler/utils"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"
)

type Store interface {
	GetRobotParameter(rid string) (*RobotParameter, error)
	CreateRobotParameter(para *RobotParameter) error
	UpdateRobotParameter(rid string, para *RobotParameter) error // 不能更改addr和key
	GetPoolById(id int) (*PoolInfo, error)
	GetPoolByAddr(addr, token string) (*PoolInfo, error)
	LoadAllRobot() ([]string, error)
}

type StoreImpl struct {
	db *sql.DB
	l  *utils.Logger
	e  Encoder
}

func NewStore(lc fx.Lifecycle, log mylog.Logging, cfg config.ConfigI) Store {
	s := &StoreImpl{
		l: log.GetLogger("db"),
	}

	lc.Append(fx.Hook{
		// app.start调用
		OnStart: func(ctx context.Context) error {
			// 这里不能阻塞
			var err error
			s.db, err = sql.Open("mysql", cfg.GetElem("db").(string))
			return err
		},
		// app.stop调用，收到中断信号的时候调用app.stop
		OnStop: func(ctx context.Context) error {
			return s.db.Close()
		},
	})

	return s
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

/*
CREATE TABLE IF NOT EXISTS `robot` (
    `id` INT(11) AUTO_INCREMENT,
    `rid` VARCHAR(20) NOT NULL COMMENT 'robot id',
    `pool_id` INT(11) NOT NULL COMMENT 'pool id',
    `start_num` INT(11) NOT NULL COMMENT '第n次开始下注',
    `num_of_bets` INT(11) NOT NULL COMMENT '最多连续下注n手',
    `addr` VARCHAR(64) NOT NULL COMMENT 'bet addr',
    `key` VARCHAR(128) NOT NULL COMMENT 'encrypto private key',
    `odd_chips` VARCHAR(128) NOT NULL COMMENT 'eg:21-41-81-161',
    `even_chips` VARCHAR(128) NOT NULL COMMENT 'eg:20-40-80-160',
    `take_profit` INT(11) NOT NULL COMMENT '止盈',
    `stop_loss` INT(11) NOT NULL COMMENT '止损',
    `state` INT(2) NOT NULL DEFAULT 0 COMMENT '0:stop/1:run',
    `ts` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP comment 'create time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `rid` (`rid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
*/

func (s *StoreImpl) GetRobotParameter(rid string) (*RobotParameter, error) {
	var r RobotParameter
	r.Rid = rid

	var oc, ec string
	row := s.db.QueryRow("SELECT `pool_id`,`start_num`,`num_of_bets`,`addr`,`key`,`odd_chips`,`even_chips`,`take_profit`,`stop_loss`,`state` FROM `robot` WHERE `rid`=?", rid)

	err := row.Scan(&r.PoolId, &r.StartNum, &r.NumOfBets, &r.Addr, &r.Key, &oc, &ec, &r.TP, &r.SL, &r.State)
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, utils.LineNo())
	}

	r.OddChips, err = s.e.IntArray(oc)
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, utils.LineNo())
	}

	r.EvenChips, err = s.e.IntArray(ec)
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, utils.LineNo())
	}
	return &r, nil
}

func (s *StoreImpl) CreateRobotParameter(r *RobotParameter) error {
	sql := "INSERT INTO `robot` (`rid`,`pool_id`,`start_num`,`num_of_bets`,`addr`,`key`,`odd_chips`,`even_chips`,`take_profit`,`stop_loss`) VALUES (?,?,?,?,?,?,?,?,?,?)"
	oc := s.e.String(r.OddChips)
	ec := s.e.String(r.EvenChips)
	_, err := s.db.Exec(sql, r.Rid, r.PoolId, r.StartNum, r.NumOfBets, r.Addr, r.Key, oc, ec, r.TP, r.SL)
	if err != nil {
		return fmt.Errorf("%w%s", err, utils.LineNo())
	}

	return nil
}

func (s *StoreImpl) UpdateRobotParameter(rid string, para *RobotParameter) error {
	sql := "UPDATE `robot` SET `pool_id`=?,`start_num`=?,`num_of_bets`=?,`addr`=?,`key`=?,`odd_chips`=?,`even_chips`=?,`take_profit`=?,`stop_loss`=?,`state`=? WHERE `rid` = ?"
	oc := s.e.String(para.OddChips)
	ec := s.e.String(para.EvenChips)
	_, err := s.db.Exec(sql, para.PoolId, para.StartNum, para.NumOfBets, para.Addr, para.Key, oc, ec, para.TP, para.SL, para.State, rid)
	if err != nil {
		return fmt.Errorf("%w%s", err, utils.LineNo())
	}

	return nil
}

func (s *StoreImpl) LoadAllRobot() ([]string, error) {
	sql := "SELECT `rid` FROM `robot` WHERE `state`=1"
	rs := make([]string, 0, 16)

	rows, err := s.db.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, utils.LineNo())
	}

	var tmp string
	for rows.Next() {
		_ = rows.Scan(&tmp)
		rs = append(rs, tmp)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%w%s", err, utils.LineNo())
	}

	return rs, nil

}

type PoolInfo struct {
	Id        int
	Kind      int
	Addr      string
	Refund    string
	Token     string
	MinAmount int
	MaxAmount int
	State     int
}

func (s *StoreImpl) GetPoolById(id int) (*PoolInfo, error) {
	var p PoolInfo
	sql := "SELECT `id`,`kind`,`addr`,`refund`,`token`,`min_amount`,`max_amount`,`state` FROM `pool` WHERE `id` = ?"

	row := s.db.QueryRow(sql, id)
	err := row.Scan(&p.Id, &p.Kind, &p.Addr, &p.Refund, &p.Token, &p.MinAmount, &p.MaxAmount, &p.State)
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, utils.LineNo())
	}
	return &p, nil
}

func (s *StoreImpl) GetPoolByAddr(addr, token string) (*PoolInfo, error) {
	var p PoolInfo
	sql := "SELECT `id`,`kind`,`addr`,`refund`,`token`,`min_amount`,`max_amount`,`state` FROM `pool` WHERE `addr`=? AND `token`=?"
	row := s.db.QueryRow(sql, addr, token)
	err := row.Scan(&p.Id, &p.Kind, &p.Addr, &p.Refund, &p.Token, &p.MinAmount, &p.MaxAmount, &p.State)
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, utils.LineNo())
	}
	return &p, nil
}

var StoreModule = fx.Options(fx.Provide(NewStore))
