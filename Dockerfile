FROM debian:jessie
COPY bin/gearman-load-logger /usr/local/bin/gearman-load-logger
COPY ./kvconfig.yml /usr/local/bin/kvconfig.yml
CMD ["/usr/local/bin/gearman-load-logger"]
