FROM alpine:latest
MAINTAINER Frederic Delbos <fred.delbos@gmail.com>

RUN apk --no-cache add ca-certificates

COPY piped /bin/piped
EXPOSE 7890

ENTRYPOINT ["/bin/piped"]
