FROM golang:latest as builder
MAINTAINER daniel840829 "s102033114@gapp.nthu.edu.tw"
COPY . $GOPATH/src/github.com/daniel840829/gameServer
WORKDIR $GOPATH/src/github.com/daniel840829/gameServer
RUN set -x && go get github.com/golang/dep/cmd/dep && dep ensure -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o /gameServer
CMD ["/gameServer"]
EXPOSE 3000 8080 50051

FROM alpine
COPY --from=builder /gameServer .
COPY ./cluster ./cluster
EXPOSE 3000 8080 50051
# ENTRYPOINT [ "/bin/bash" ]
CMD ["./gameServer"]


