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

#### table 结构

| name       | type         |  notes      |
| --------   | -----        | :----:      |
| id         | int(11)      | atom        |
| rid        | varchar(20)  | robot id    |
| user       | varchar(32)  | belong to   |
| token      | varchar(32)  | 币种        |
| count      | int(11)      | 8           |
| addr       | varchar(64)  | 下注地址    |
| key        | varchar(128) | 私钥        |
| take_profit| int(11)      | 止盈        |
| stop_loss  | int(11)      | 止损        |
| odd_bets   | varchar(128) | 21-41-81-161|
| even_bets  | varchar(128) | 20-40-80-160|
| state      | int(2)       | 0:stop,1:run|
| ts         | int(11)      | create time |
| md5        | varchar(64)  | md5(user+token+count+odd_bets+even_bets)|

```
CREATE TABLE IF NOT EXISTS `robot` (
    `id` INT(11) AUTO_INCREMENT,
    `rid` VARCHAR(20) NOT NULL COMMENT 'robot id',
    `user` VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'belone to',
    `token` VARCHAR(32) NOT NULL COMMENT 'TRX/USDT',
    `count` INT(11) NOT NULL COMMENT '连续n次奇/偶后开始下注',
    `addr` VARCHAR(64) NOT NULL COMMENT 'bet addr',
    `key` VARCHAR(128) NOT NULL COMMENT 'encrypto private key',
    `odd_bets` VARCHAR(128) NOT NULL COMMENT 'eg:21-41-81-161',
    `even_bets` VARCHAR(128) NOT NULL COMMENT 'eg:20-40-80-160',
    `take_profit` INT(11) NOT NULL COMMENT '止盈',
    `stop_loss` INT(11) NOT NULL COMMENT '止损',
    `state` INT(2) NOT NULL DEFAULT 0 COMMENT '0:stop/1:run',
    `ts` INT(11) NOT NULL COMMENT 'create time',
    `md5` VARCHAR(64) NOT NULL COMMENT 'md5(user+token+count+odd_bets+even_bets)',
    PRIMARY KEY (`id`),
    UNIQUE KEY `md5` (`md5`),
    KEY (`rid`),
    KEY (`user`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## API
#### 创建一个机器人
* POST /robot/create
* post body:
```
{
    "token":        "TRX/USDT",
    "count":        n,           //连续n次奇数/偶数后开始下注
    "odd_bets":     [21,41,81,161...], //奇数赌注
    "even_bets":    [20,40,80,160...], //偶数赌注
    "take_profit":  50000, //止盈
    "stop_loss":    10000, //止损
    "ts":           1675052150 //创建时间
}
```


* response:
```
{
    "code":         200,  (200: 成功)  (xxx: 失败)
    "msg":          "success/error",
    "version":      "2.0.0"
    "id":           "xxxxxxx"  // 机器人唯一ID
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
    "id":           "xxxxxxx"  // 机器人唯一ID
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
    "id":           "xxxxxxx"  // 机器人唯一ID
}
```

#### 获取地址

* GET /robot/address?id="xxxxxxx"
* response:

```
{
"code":         200,  (200: 成功)  (xxx: 失败)
"msg":          "success/error",
"version":      "2.0.0"
"address":      "41xxxx"  // 机器人地址
}
```

#### 获取余额（跳到波场浏览器看）
* GET /robot/balance?id="xxxxxxx"

暂未实现


#### 获取交易记录（跳到波场浏览器看）
* GET /robot/log?id="xxxxxxx"&from=1675051000&to=1675058000

暂未实现