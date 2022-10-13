build:
	docker build -t scheduler .

push:
	docker tag scheduler registry.cn-zhangjiakou.aliyuncs.com/eininst/scheduler:v1
	docker push registry.cn-zhangjiakou.aliyuncs.com/eininst/scheduler:v1

run:
	docker run  -v /Users/wangziqing/go/scheduler/configs/config.yaml:/config.yaml \
		-e profile=dev -p "3000:3000" --net=bridge \
		registry.cn-zhangjiakou.aliyuncs.com/eininst/scheduler:v1

clean:
	yes | docker system prune

.PHONY: build push run clean