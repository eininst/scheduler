# Scheduler

[![Build Status](https://travis-ci.org/ivpusic/grpool.svg?branch=master)](https://github.com/infinitasx/easi-go-aws)

## ⚙ Installation

```docker
docker pull registry.cn-zhangjiakou.aliyuncs.com/eininst/scheduler:v1
```

## ⚡ Quickstart

```docker
docker run  -v /xxx/xxx.yaml:/config.yaml \
    -e profile=dev -p "3000:3000" --net=bridge \
    registry.cn-zhangjiakou.aliyuncs.com/eininst/scheduler:v1
```
## Config yaml

```text
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

## License

*MIT*