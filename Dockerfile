FROM alpine:3.4

RUN apk --no-cache add go git &&\
 mkdir /go &&\
 export GOPATH=/go &&\
 go get -v github.com/fern4lvarez/piladb/pilad &&\
 go build -o /pilad -ldflags "-s -w" github.com/fern4lvarez/piladb/pilad &&\
 rm -rf /go &&\
 apk del go git

ENTRYPOINT ["/pilad"]