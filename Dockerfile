FROM ubuntu:16.04

MAINTAINER Nate Sweet <nathanjsweet@gmail.com>

CMD ["/bin/healthcheck-app"]

COPY healthcheck-app /bin/healthcheck-app

