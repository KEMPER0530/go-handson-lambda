FROM golang:1.13.4-alpine3.10
MAINTAINER kemper0530

ENV GOPATH /go
ENV PATH=$PATH:$GOPATH/src

RUN apk update && \
    apk upgrade && \
    apk add vim && \
    apk add git && \
    apk add build-base

# Install our dependencies
RUN go get github.com/gin-gonic/gin
RUN go get github.com/jinzhu/gorm
RUN go get github.com/go-sql-driver/mysql
RUN go get golang.org/x/tools/cmd/goimports
RUN go get github.com/joho/godotenv
RUN go get golang.org/x/crypto/bcrypt
RUN go get github.com/google/uuid
RUN go get github.com/aws/aws-sdk-go/aws
RUN go get github.com/aws/aws-sdk-go/aws/session
RUN go get github.com/aws/aws-sdk-go/service/ses
RUN go get github.com/bamzi/jobrunner
RUN go get github.com/k-washi/jwt-decode/jwtdecode
RUN go get firebase.google.com/go
RUN go get github.com/gin-contrib/cors
RUN go get github.com/aws/aws-lambda-go/events
RUN go get github.com/aws/aws-lambda-go/lambda
RUN go get github.com/awslabs/aws-lambda-go-api-proxy/gin

#RUN go get -u github.com/golang/protobuf/proto
#RUN go get -u google.golang.org/grpc
#RUN go get -u github.com/grpc-ecosystem/go-grpc-middleware
#RUN go get -u github.com/urfave/cli
#RUN go get -u google.golang.org/grpc/credentials
#RUN go get -u google.golang.org/grpc/reflection
#RUN go get -u github.com/xo/xo
#RUN go get -u github.com/jmoiron/sqlx
#RUN go get -u github.com/Masterminds/squirrel

RUN apk add tzdata
ENV TZ=Asia/Tokyo

# 以下、Docker run 用の設定
#ENV GO_ENV=production
ENV PATH=$PATH:$GOPATH/src/github.com/kemper0530/go-handson-lambda

WORKDIR $GOPATH/src/github.com/kemper0530/go-handson-lambda

COPY  /config $GOPATH/src/github.com/kemper0530/go-handson-lambda/config
COPY  /keys $GOPATH/src/github.com/kemper0530/go-handson-lambda/keys
COPY  /common $GOPATH/src/github.com/kemper0530/go-handson-lambda/common
COPY  /controllers $GOPATH/src/github.com/kemper0530/go-handson-lambda/controllers
COPY  /models $GOPATH/src/github.com/kemper0530/go-handson-lambda/models
COPY  /static $GOPATH/src/github.com/kemper0530/go-handson-lambda/static
COPY  /templates $GOPATH/src/github.com/kemper0530/go-handson-lambda/templates
COPY  /serve_test.go $GOPATH/src/github.com/kemper0530/go-handson-lambda
COPY  /serve.go $GOPATH/src/github.com/kemper0530/go-handson-lambda

RUN go mod init github.com/kemper0530/go-handson-lambda
RUN GOOS=linux go build -o go-handson-lambda
RUN chmod -R 777 $GOPATH/src/github.com/kemper0530/go-handson-lambda

ENTRYPOINT ["/go/src/github.com/kemper0530/go-handson-lambda/go-handson-lambda"]

# Expose default port (8080)
#EXPOSE 8080
#EXPOSE 8090
