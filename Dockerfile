#源镜像
FROM golang:latest
#作者
MAINTAINER daniel840829 "s102033114@gapp.nthu.edu.tw"
#设置工作目录
WORKDIR $GOPATH/src/github.com/daniel840829/gameServer
#将服务器的go工程代码加入到docker容器中
ADD . $GOPATH/src/github.com/daniel840829/gameServer
#go构建可执行文件
RUN go get github.com/daniel840829/gameServer
RUN go install
#暴露端口
EXPOSE 3000 8080 50051
#最终运行docker的命令
ENTRYPOINT  ["/go/bin/gameServer"]
