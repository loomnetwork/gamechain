FROM golang:latest
ADD . /gamechain
WORKDIR /gamechain
RUN make deps
RUN make gamechain-logger
CMD [ "bin/gamechain-logger", "${WEBSOCKET_URL}"]