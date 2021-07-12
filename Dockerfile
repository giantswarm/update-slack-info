FROM golang:1.14-alpine3.11@sha256:6578dc0c1bde86ccef90e23da3cdaa77fe9208d23c1bb31d942c8b663a519fa5 AS builder

LABEL "com.github.actions.icon"="bell"
LABEL "com.github.actions.color"="yellow"
LABEL "com.github.actions.name"="Update Slack Info"
LABEL "com.github.actions.description"="This action will update Slack Information"
LABEL "org.opencontainers.image.source"="https://github.com/giantswarm/update-slack-info"

WORKDIR ${GOPATH}/src/github.com/giantswarm/update-slack-info
COPY main.go ${GOPATH}/src/github.com/giantswarm/update-slack-info

ENV CGO_ENABLED 0
ENV GOOS linux

RUN go get -v ./...
RUN go build -a -installsuffix cgo -ldflags '-w  -extldflags "-static"' -o /go/bin/update-slack-info .

# alpine:latest at 2020-01-18T01:19:37.187497623Z
FROM alpine@sha256:ab00606a42621fb68f2ed6ad3c88be54397f981a7b70a79db3d1172b11c4367d

COPY --from=builder /go/bin/update-slack-info /usr/bin/update-slack-info

RUN apk update \
	&& apk upgrade \
	&& apk add \
	bash \
	jq \
	ca-certificates \
	python \
	py2-pip \
	rsync && \
	pip install shyaml && \
	rm -rf /var/cache/apk/*

COPY *.sh /

RUN chmod +x /*.sh

ENTRYPOINT ["/entrypoint.sh"]
