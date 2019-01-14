FROM golang:latest
RUN make deps
RUN make gamechain-logger
CMD [ "bin/gamechain-logger", "${WEBSOCKET_URL}"]