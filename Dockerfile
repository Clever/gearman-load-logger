FROM debian:jessie
COPY bin/gearman-load-logger /usr/local/bin/gearman-load-logger
CMD ["/usr/local/bin/gearman-load-logger"]
