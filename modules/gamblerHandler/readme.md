### 创建一个机器人
* POST /robot/create
* post body:

```
{
    "token":        "TRX/USDT",
    "count":        n,           //连续n次奇数/偶数后开始下注
    "odd_bets":     [21,41,81,161...], //奇数赌注
    "enen_bets":    [20,40,80,160...], //偶数赌注
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

### 开启一个机器人
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


### 停止一个机器人
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

### 获取余额
* GET /robot/balance?id="xxxxxxx"