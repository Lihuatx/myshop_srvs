#定义一个变量
srv_name="user-srv-main"
#添加可执行权限
chmod +x ./$srv_name
#重启，如果已经在则关闭重启
if pgrep -x $srv_name > /dev/null
then
	echo "${srv_name} is running"
	#先关闭进程
	echo "shutting down ${srv_name}"
	#这里不要使用kill -9，这个是强杀，我们需要优雅退出
	if kill $(pgrep -f $srv_name)
		then
			echo "starting ${srv_name}"
			#以后台方式启动，定向到dev/null 2>&1 &
			./$srv_name > /dev/null 2>&1 &
			echo "start ${srv_name} success"
	fi
else
	echo "starting ${srv_name}"
	./$srv_name > /dev/null 2>&1 &
	echo "start ${srv_name} success"
fi

