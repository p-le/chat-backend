FROM golang:1.8

COPY ./app /go/src/app
WORKDIR /go/src/app

RUN go get ./
RUN go build
RUN go get github.com/beego/bee

EXPOSE 2712
CMD ["bee", "run"]