#!/bin/bash

# 定义服务名称数组
services=("goods-srv-main" "inventory-srv-main" "order-srv-main" "user-srv-main" "userop-srv-main")

echo "开始停止所有服务..."

# 遍历停止所有服务
for service in "${services[@]}"
do
    if pgrep -x "$service" > /dev/null
    then
        echo "正在停止 $service..."
        # 使用kill而不是kill -9，实现优雅退出
        if ps -a | grep "$service" | awk '{print $1}' | xargs kill
        then
            echo "$service 已停止"
        else
            echo "警告: $service 停止失败"
        fi
    else
        echo "$service 未运行"
    fi
done

# 等待几秒确保服务都已停止
sleep 3

# 检查服务是否完全停止
echo "检查服务状态："
for service in "${services[@]}"
do
    if pgrep -x "$service" > /dev/null
    then
        echo "警告: $service 仍在运行"
    else
        echo "$service 已完全停止"
    fi
done

echo "停止服务操作完成！"