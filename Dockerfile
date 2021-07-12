FROM quay.io/giantswarm/alpine:3.12

RUN apk --no-cache add ca-certificates

RUN mkdir -p /opt
WORKDIR /opt

RUN mkdir -p /opt
COPY ./update-slack-info /opt

ENTRYPOINT ["/opt/update-slack-info"]
