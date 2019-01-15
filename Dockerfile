FROM golang:latest
ADD bin/gamechain-logger gamechain-logger
CMD [ "./gamechain-logger", "${WEBSOCKET_URL}"]