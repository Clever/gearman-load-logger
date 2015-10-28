FROM golang:1.5
RUN apt-get update
RUN apt-get install -y curl build-essential

# copy source code
RUN mkdir -p /go/src/github.com/Clever/gearman-load-logger
ADD . /go/src/github.com/Clever/gearman-load-logger

WORKDIR /go/src/github.com/Clever/gearman-load-logger

RUN go install github.com/Clever/gearman-load-logger
RUN go build -o /usr/local/bin/gearman-load-logger github.com/Clever/gearman-load-logger

CMD ["/usr/local/bin/gearman-load-logger"]
