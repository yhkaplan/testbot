FROM alpine:3.4

RUN apk add --no-cache ca-certificates

ADD testbot testbot
RUN chmod +x testbot

CMD ["./testbot"]
