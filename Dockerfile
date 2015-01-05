FROM golang:1.3
RUN apt-get update
RUN apt-get install -y curl build-essential
RUN go get github.com/tools/godep

# Gearcmd
RUN curl -L https://github.com/Clever/gearcmd/releases/download/v0.3.1/gearcmd-v0.3.1-linux-amd64.tar.gz | tar xz -C /usr/local/bin --strip-components 1

# copy source code
RUN mkdir -p /go/src/github.com/Clever/gearman-load-logger
ADD . /go/src/github.com/Clever/gearman-load-logger

# set workdir to find saved godeps
WORKDIR /go/src/github.com/Clever/gearman-load-logger

# build source code using godep
RUN rm -rf /go/src/github.com/Clever/gearman-load-logger/Godeps/_workspace/pkg/
RUN godep go install github.com/Clever/gearman-load-logger
RUN godep go build -o /usr/local/bin/gearman-load-logger github.com/Clever/gearman-load-logger

CMD ["/go/src/github.com/Clever/gearman-load-logger/run_as_worker.sh"]
