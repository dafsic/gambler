# gambler --- TRX/USDT 自动下注机器人平台

## 逻辑
#### 机器人逻辑
1. 监听所有区块，如果区块hash值最后一个数字连续n次为奇数/偶数，则开始下注
2. 如果中了，则重新开始计数
3. 如果没中，赌注加倍，在下一个区块继续下注，直到中了为止

#### 控制逻辑
1. 接收机器人配置，写入数据库，然后开启一个机器人，状态设置为1，返回机器人id
2. 接收到停止某id的机器人，关闭线程，数据库状态设置为0
3. 程序启动时，从数据库中加载所有状态为1的机器人配置

#### 业务逻辑
1. 池子是我们自己定义的，用户无法自定义
2. 机器人可以包月或者包天收费
3. 池子方可以上池收费，给我钱我就上你的池子

## 版本1（平台提供池子）
#### robot table 结构

| name       | type         |  notes      |
| --------   | -----        | :----:      |
| id         | int(11)      | atom        |
| rid        | varchar(20)  | robot id    |
| pool_id    | int(11)      | pool id     |
| start_num  | int(11)      | 8           |
| num_of_bets| int(11)      | 最多下注次数|
| addr       | varchar(64)  | 下注地址    |
| key        | varchar(128) | 私钥        |
| take_profit| int(11)      | 止盈        |
| stop_loss  | int(11)      | 止损        |
| odd_chips  | varchar(128) | 21-41-81-161|
| even_chips | varchar(128) | 20-40-80-160|
| state      | int(2)       | 0:stop,1:run|
| ts         | timestamp    | create time |

```
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
```

#### pool table 结构

| name       | type         |  notes      |
| --------   | -----        | :----:      |
| id         | int(11)      | atom        |
| kind       | int(11)      | robot kind  |
| addr       | varchar(64)  | 池地址      |
| refund     | varchar(64)  | 回款地址    |
| token      | varchar(32)  | 币种        |
| min_amount | int(11)      | 最小下注    |
| max_amount | int(11)      | 最大下注    |
| state      | int(2)       | 0:stop,1:run|
| ts         | timestamp    | create time |

```
CREATE TABLE IF NOT EXISTS `pool` (
    `id` INT(11) AUTO_INCREMENT,
    `addr` VARCHAR(64) NOT NULL COMMENT 'pool addr',
    `kind` int(11) NOT NULL DEFAULT 1 COMMENT '规则',
    `refund` VARCHAR(64) NOT NULL COMMENT 'refund addr',
    `token` VARCHAR(32) NOT NULL COMMENT 'TRX/USDT',
    `min_amount` INT(11) NOT NUL COMMENT '最小下注',
    `max_amount` INT(11) NOT NUL COMMENT '最小下注',
    `state` INT(2) NOT NULL DEFAULT 0 COMMENT '0:invalid/1:valid',
    `ts` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP comment 'create time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `addr` (`addr`,`token`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## API
#### 创建一个机器人(不会自动开启)
* POST /robot/create
* post body:
```
{
    "pool":          "41xxxxxx",  //池地址
    "token":         "TRX/USDT",  //币种
    "start_num":     m,           //第m次开始下注
    "num_of_bets":   n,           //每轮最多下n次
    "odd_chips":     [21,41,81,161...], //奇数筹码
    "even_chips":    [20,40,80,160...], //偶数筹码
    "take_profit":  50000, //止盈
    "stop_loss":    10000, //止损
}
```


* response:
```
{
    "code":         200,  (200: 成功)  (xxx: 失败)
    "msg":          "success/error",
    "version":      "2.0.0"
    "id":           "xxxxxxx"  // 机器人唯一ID
    "addr":         "41xxxxx"  // 机器人地址
}
```

#### 开启一个机器人
* GET /robot/run?id="xxxxxxx"
* response:
```
{
    "code":         200,  (200: 成功)  (xxx: 失败)
    "msg":          "success/error",
    "version":      "2.0.0"
}
```


#### 停止一个机器人
* GET /robot/stop?id="xxxxxxx"
* response:
```
{
    "code":         200,  (200: 成功)  (xxx: 失败)
    "msg":          "success/error",
    "version":      "2.0.0"
}
```

#### 更新机器人配置（下次启动后生效）

* POST /robot/update?id="xxxxxxx"
* post body:
```
{
    "pool":          "41xxxxxx",  //池地址
    "token":         "TRX/USDT",  //币种
    "start_num":     m,           //连续m次奇数/偶数后开始下注
    "stop_num":      n,           //最多连续下n次
    "odd_chips":     [21,41,81,161...], //奇数筹码
    "even_chips":    [20,40,80,160...], //偶数筹码
    "take_profit":  50000, //止盈
    "stop_loss":    10000, //止损
}
```

* response:
```
{
"code":         200,  (200: 成功)  (xxx: 失败)
"msg":          "success/error",
"version":      "2.0.0"
}
```

#### 获取余额(只支持trx/usdt)
* GET /robot/balance?addr="41xxxxxxx"&token="trx"
* response:
```
{
"code":         200,  (200: 成功)  (xxx: 失败)
"msg":          "success/error",
"version":      "2.0.0"
"balance":      100.62
}
```


#### 获取交易记录（跳到波场浏览器看）
* GET /robot/log?id="xxxxxxx"&from=1675051000&to=1675058000

暂未实现


## 版本2（用户自己提供池子）
#### robot2 table 结构

| name       | type         |  notes      |
| --------   | -----        | :----:      |
| id         | int(11)      | atom        |
| rid        | varchar(20)  | robot id    |
| pool       | varchar(64)  | 池地址      |
| refund     | varchar(64)  | 回款地址    |
| token      | varchar(32)  | 币种        |
| min_amount | int(11)      | 最小下注    |
| max_amount | int(11)      | 最大下注    |
| start_num  | int(11)      | 8           |
| num_of_bets| int(11)      | 最多下注次数|
| addr       | varchar(64)  | 下注地址    |
| key        | varchar(128) | 私钥        |
| take_profit| int(11)      | 止盈        |
| stop_loss  | int(11)      | 止损        |
| odd_chips  | varchar(128) | 21-41-81-161|
| even_chips | varchar(128) | 20-40-80-160|
| state      | int(2)       | 0:stop,1:run|
| ts         | timestamp    | create time |

```
CREATE TABLE IF NOT EXISTS `robot2` (
    `id` INT(11) AUTO_INCREMENT,
    `rid` VARCHAR(20) NOT NULL COMMENT 'robot id',
    `pool` VARCHAR(64) NOT NULL COMMENT 'pool addr',
    `refund` VARCHAR(64) NOT NULL COMMENT 'refund addr',
    `token` VARCHAR(32) NOT NULL COMMENT 'TRX/USDT',
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
```