#源镜像
FROM golang:latest as builder
#作者
MAINTAINER daniel840829 "s102033114@gapp.nthu.edu.tw"
#设置工作目录
COPY . $GOPATH/src/github.com/daniel840829/gameServer
WORKDIR $GOPATH/src/github.com/daniel840829/gameServer
#将服务器的go工程代码加入到docker容器中

RUN set -x && go get github.com/golang/dep/cmd/dep && dep ensure -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o /gameServer


CMD ["/gameServer"]

#暴露端口
EXPOSE 3000 8080 50051
# 最终运行docker的命令
# ENTRYPOINT  ["/go/bin/gameServer"]
# ---stage 2
FROM scratch

COPY --from=builder /gameServer .


EXPOSE 3000 8080 50051

CMD ["./gameServer"]
