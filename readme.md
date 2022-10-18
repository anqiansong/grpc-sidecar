# grpc-sidecar

本工程仅供学习使用,目前仅实现了简单的拦截，服务鉴定，限流等功能，这些功能只起演示作用，其他的功能可实现 Filter 接口来实现

# 依赖

1. 本地 etcd-server
2. 本地 redis-server

#  本地启动顺序（非 k8s 环境）

1. 启动 proxy -> main.go
2. 启动 service server -> example/server/main.go
3. 启动 client -> example/client/main.go
4. 启动 cp 程序模拟配置下发 -> example/cp/main.go

