###########################################################################
# This file moved to root so that it can get full context of source files #
###########################################################################
##########################
#      Build Stage       #
##########################
FROM golang:1.10-alpine AS base

RUN apk update && \
    apk add openssh gcc libc-dev cyrus-sasl-dev linux-headers python-dev automake autoconf inotify-tools py2-pip openssl-dev git make tzdata
RUN apk add --no-cache mongodb &&\
    mkdir -p /data/db


RUN go get -u github.com/golang/dep/cmd/dep
RUN go get github.com/pilu/fresh
RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

COPY ./ $GOPATH/src/github.com/lucasmoten/project-2502/services/dice
WORKDIR $GOPATH/src/github.com/lucasmoten/project-2502/services/dice

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.0.0/wait /wait
RUN chmod +x /wait

WORKDIR $GOPATH/src/github.com/lucasmoten/project-2502/services/dice/

CMD fresh cmd/server/main.go --config settings.toml

#################################
#  Build Binary 4 docker image  #
#################################

FROM base as bin
# build binary
RUN go build -v -o cmd/server/dice-bin ./cmd/server


##############################
#  Deployable Docker Image#  #
##############################
FROM alpine:3.7

# copy over needed materials (bin, settings, certs, entry, docs)
COPY --from=bin /go/src/github.com/lucasmoten/project-2502/services/dice/cmd/server/dice-bin /service/dice
COPY ./docker/build/settings.toml /service/settings.toml
COPY ./docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh && chmod +x /service/dice

EXPOSE 10000
EXPOSE 10001
EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]