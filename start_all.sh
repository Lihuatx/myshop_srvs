#!/bin/bash

echo "启动 goods-srv 服务..."
cd ./goods_srv/target
./start.sh
sleep 1

echo "启动 inventory-srv 服务..."
cd ../../inventory_srv/target
./start.sh
sleep 1

echo "启动 order-srv 服务..."
cd ../../order_srv/target
./start.sh
sleep 1

echo "启动 user-srv 服务..."
cd ../../user_srv/target
./start.sh
sleep 1

echo "启动 userop-srv 服务..."
cd ../../userop_srv/target
./start.sh

echo "所有服务启动完成！"

# 等待所有服务启动
sleep 5

# 检查服务状态
echo "检查服务状态："
for service in "goods-srv-main" "inventory-srv-main" "order-srv-main" "user-srv-main" "userop-srv-main"
do
    if pgrep -x "$service" > /dev/null
    then
        echo "$service 运行正常"
    else
        echo "警告: $service 可能未正常启动"
    fi
done