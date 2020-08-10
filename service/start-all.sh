#!/bin/bash

# 检查service进程
check_process() {
    sleep 1
    res=`ps aux | grep -v grep | grep "service/bin" | grep $1`
    if [[ $res != '' ]]; then
        echo -e "\033[32m 已启动: \033[0m" "$1"
        return 1
    else
        echo -e "\033[31m 启动失败: \033[0m" "$1"
        return 0
    fi
}

# 编译service可执行文件
build_service() {
    echo -e "\033[32m 开始编译: \033[0m service/bin/$resbin"
    go build -o service/bin/$1 service/$1/main.go
    resbin=`ls service/bin/ | grep $1`
    echo -e "\033[32m 编译完成: \033[0m service/bin/$resbin"
}

# 启动service
run_service() {
    nohup ./service/bin/$1 --registry=consul >> $logpath/$1.log 2>&1 &
    sleep 1
    check_process $1
}

ignoreGOPATH=$1
if [[ $ignoreGOPATH == '-ignoreGOPATH' ]];
then
	echo -e "\033[31m WARNING: \033[0m \$GOPATH ignored..."
else
	export LD_LIBRARY_PATH=/usr/local/lib
	export PATH=$GOPATH/bin:$PATH
	# 切换到工程根目录
	cd $GOPATH/src/filestore-server
	#cd /data/go/work/src/filestore-server
fi

# 创建运行日志目录
logpath=/data/log/filestore-server
mkdir -p $logpath


# 微服务可以用supervisor做进程管理工具；
# 或者也可以通过docker/k8s进行部署

services="
dbproxy
upload
download
transfer
account
apigw
"

# 打包静态资源文件
res=$(rm -rf assets && mkdir assets -p && go-bindata-assetfs -pkg assets -o ./assets/asset.go static/...)
if [[ $res == "" ]]; then
    echo '静态资源文件打包完成，输出目录: assets/'
else
    exit 1
fi

# 执行编译service
res=$(mkdir -p service/bin/ && rm -f service/bin/*)
if [[ $res != "" ]]; then
    exit 1
fi

for service in $services
do
    build_service $service
done

# 执行启动service
for service in $services
do
    run_service $service
done

echo '微服务启动完毕.'

