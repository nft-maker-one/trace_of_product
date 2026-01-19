#!/bin/bash

# 启动区块链节点的脚本
echo "启动区块链节点网络..."

# 进入区块链目录
cd trace_of_product

# 构建并启动区块链节点和MySQL
docker-compose up -d --build

echo "区块链节点网络启动完成!"
echo "节点访问端口:"
echo "- 节点1: http://localhost:8081"
echo "- 节点2: http://localhost:8082"
echo "- 节点3: http://localhost:8083"
echo "- 节点4: http://localhost:8084"
echo "- MySQL: localhost:3306"

# 等待服务启动
echo "等待服务启动..."
sleep 10

# 检查服务状态
docker-compose ps