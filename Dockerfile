FROM alpine:edge
MAINTAINER Olexander Simonov <oleksandr@amoniac.eu>

RUN \
    apk add --no-cache --update ca-certificates s6 \
    && echo 'hosts: files mdns4_minimal [NOTFOUND=return] dns mdns4' >> /etc/nsswitch.conf \
    && rm -rf /var/cache/apk/*

COPY ./docker /
RUN  chmod a+x /etc/s6/.s6-svscan/finish /etc/s6/*/run /etc/s6/*/finish
COPY ./bin/slackopher /slackopher

EXPOSE 3030

CMD ["/bin/s6-svscan", "/etc/s6"]
