FROM frolvlad/alpine-glibc:alpine-3.7

ADD bin/gamechain-logger gamechain-logger

CMD [ "./gamechain-logger", "${WEBSOCKET_URL}"]
