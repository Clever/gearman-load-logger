FROM golang:1.3
RUN apt-get update
RUN apt-get install -y wget build-essential

# Gearcmd
RUN wget https://github.com/Clever/gearcmd/releases/download/v0.3.1/gearcmd-v0.3.1-linux-amd64.tar.gz
RUN tar -xvf gearcmd-v0.3.1-linux-amd64.tar.gz
RUN cp gearcmd-v0.3.1-linux-amd64/gearcmd /usr/local/bin/

# Set up worker
RUN mkdir -p /go/src/github.com/Clever/gearman-load-logger
ADD . /go/src/github.com/Clever/gearman-load-logger/
RUN GOPATH=/go go get -d github.com/Clever/gearman-load-logger
RUN GOPATH=/go go build -o /usr/local/bin/gearman-load-logger github.com/Clever/gearman-load-logger
CMD ["/go/src/github.com/Clever/gearman-load-logger/run_as_worker.sh"]
