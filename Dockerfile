FROM golang:1.3
RUN apt-get update
RUN apt-get install -y curl build-essential
RUN go get github.com/tools/godep

# copy source code
RUN mkdir -p /go/src/github.com/Clever/gearman-load-logger
ADD . /go/src/github.com/Clever/gearman-load-logger

# set workdir to find saved godeps
WORKDIR /go/src/github.com/Clever/gearman-load-logger

# build source code using godep
RUN rm -rf /go/src/github.com/Clever/gearman-load-logger/Godeps/_workspace/pkg/
RUN godep go install github.com/Clever/gearman-load-logger
RUN godep go build -o /usr/local/bin/gearman-load-logger github.com/Clever/gearman-load-logger

CMD ["/usr/local/bin/gearman-load-logger"]
