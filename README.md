# Scheduler

[![Build Status](https://travis-ci.org/ivpusic/grpool.svg?branch=master)](https://github.com/infinitasx/easi-go-aws)

## Installation
```text
go get -u github.com/eininst/scheduler
```

## ⚡ Quickstart
```go
package main

import "github.com/eininst/scheduler"

func main() {
    app := scheduler.New("./configs/config.yaml")
    app.Listen()
}
```
> runing...
```text
2022/10/13 17:59:26 [INFO] config.go:24 profile is: dev
2022/10/13 17:59:26 [DEBUG] rdb.go:47 Connected to Redis server...  addr=r-8vbbqc8wfszjqhvnvnpd.redis.zhangbei.rds.aliyuncs.com:6379 db=0 poolSize=100
2022/10/13 17:59:27 Connected to Mysql server...
2022/10/13 17:59:27 [INFO] [RS] Stream "task_run:9e45874e-4310-4067-b60c-9f099ba039c6" working... # BlockTime=15s MaxRetries=3 ReadCount=20 Timeout=5m0s Work=1024
2022/10/13 17:59:27 [INFO] [RS] Stream "task_stop:4ad63215-b771-4ff9-ac27-46a8a3189ee5" working... # BlockTime=15s MaxRetries=3 ReadCount=20 Timeout=5m0s Work=1024
2022/10/13 17:59:27 [INFO] [RS] Stream "cron_task_log" working... # BlockTime=15s MaxRetries=3 ReadCount=20 Timeout=5m0s Work=5
2022/10/13 17:59:27 [INFO] [RS] Stream "cron_task_alarm" working... # BlockTime=15s MaxRetries=3 ReadCount=20 Timeout=5m0s Work=5

 ┌───────────────────────────────────────────────────┐ 
 │                   Fiber v2.37.1                   │ 
 │               http://127.0.0.1:3000               │ 
 │       (bound on host 0.0.0.0 and port 3000)       │ 
 │                                                   │ 
 │ Handlers ............ 48  Processes ........... 1 │ 
 │ Prefork ....... Disabled  PID ............. 28765 │ 
 └───────────────────────────────────────────────────┘ 

```

## Case
> visit http://localhost:3000

<img alt="Redoc logo" src="https://nft-cj2533.oss-cn-zhangjiakou.aliyuncs.com/3.png"  width="920px" height="520px"/>



## ⚙ Config

```yaml
tablePrefix: "scheduler_"
secretKey: "xxxxxxxxxxxxxxxxxxxxxxxx"
port: 3000

web:
  title: Scheduler
  desc: 简单，开箱即用的定时任务平台
  logo:
  avatar:

log:
  retain: 10
  interval: 30
  work: 10


mail:
  host: smtp.qq.com
  port: 465
  username: eininst@qq.com
  password: xxxxxxxxxxxxxxxx
  work: 5

  cc:
    - eininst@qq.com
    - eininst@aliyun.com

---
profile: dev

redis:
  addr: localhost:6379
  db: 0
  poolSize: 100
  minIdleCount: 20

mysql:
  dsn: xx:xxxxx@tcp(xxxxxx:3306)/nft?charset=utf8mb4&parseTime=True&loc=Local
  maxIdleCount: 32
  maxOpenCount: 128
  maxLifetime: 7200
```

## For Docker

```docker
docker pull registry.cn-zhangjiakou.aliyuncs.com/eininst/scheduler:v1
```

## Run in docker

```docker
docker run \
    -v /xxx/xxx.yaml:/config.yaml \
    -e profile=dev \
    -p "3000:3000" \
    --log-opt max-size=1024m \
    --log-opt max-file=3 \
    --net=bridge \
    registry.cn-zhangjiakou.aliyuncs.com/eininst/scheduler:v1
```

> See [examples](/examples)
> 
## License

*MIT*