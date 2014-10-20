FROM golang
MAINTAINER astaxie xiemengjun@gmail.com

RUN go get github.com/astaxie/beego

RUN go get github.com/beego/beeweb

WORKDIR /go/src/github.com/beego/beeweb

RUN go build

EXPOSE 8080

CMD ["./beeweb"]
