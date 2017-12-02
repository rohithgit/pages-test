The kafka package can be imported and consumed from other golang packages.


To test locally:
The prerequisites for testing this project locally are:
- golang is installed (if you have not already: https://golang.org/dl/)
- gopath is defined (if you have not already: https://golang.org/doc/code.html#GOPATH)


Then this "kafka" folder should be copied into your $GOPATH.

Also we need a few libraries to run the code and test. Please execute these:
go get -u github.com/Sirupsen/logrus
go get -u github.com/stretchr/testify
go get -u github.com/linkedin/goavro
go get -u github.com/Shopify/sarama


Then to run the tests execute:
go test

Then to view code coverage execute:
go test -cover

To generate junit compliant test report execute:
go test -v | go-junit-report > report.xml

To generate coverage report which jenkins can understand:
gocov test | gocov-xml > coverage.xml