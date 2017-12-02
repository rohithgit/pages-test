# This image "golang:1.6.0-alpine" is built from the alpine docker image.
# It has go 1.6 installed, $PATH and $GOPATH configured.
# We are just leveraging the setup which is part of the image to put our code under $GOPATH/src/
FROM golang:1.6.0-alpine
RUN apk add --update git
ADD . $GOPATH/src/etcd
WORKDIR $GOPATH/src/etcd
RUN go get -u github.com/coreos/etcd/client
RUN go get -u github.com/Sirupsen/logrus
RUN go get -u golang.org/x/net/context
RUN go get -u github.com/stretchr/testify
RUN go get -u github.com/jstemmer/go-junit-report
RUN go get -u github.com/axw/gocov/...
RUN go get -u github.com/AlekSi/gocov-xml
