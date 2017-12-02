FROM golang:1.6.0-alpine
RUN apk add --update git
ADD . $GOPATH/src/logging
WORKDIR $GOPATH/src/logging
RUN go get -u github.com/Sirupsen/logrus
RUN go get -u github.com/stretchr/testify
RUN go get -u github.com/jstemmer/go-junit-report
RUN go get -u github.com/axw/gocov/...
RUN go get -u github.com/AlekSi/gocov-xml
