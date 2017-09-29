FROM busybox:ubuntu-14.04

MAINTAINER Nate Sweet <nathanjsweet@gmail.com>

RUN echo "nobody:x:1:1:nobody:/:/bin/sh" >> /etc/passwd

RUN echo "nobody:x:1:" >> /etc/group

USER nobody

CMD ["/bin/healthcheck-app"]

COPY ./lib /lib

COPY healthcheck-app /bin/healthcheck-app

